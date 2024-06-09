package scheduler

import (
	"os"
	"syscall"
	"time"

	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"

	"github.com/go-co-op/gocron/v2"
	logr "github.com/sirupsen/logrus"
)

type Scheduler interface {
	Run(label string)
}

type RuntimeJob func(s *SchedulerProcess, job *types.Job, monitor *types.Monitor, params map[string]string)

type SchedulerProcess struct {
	Label    string
	Db       storage.Provider
	Log      *logr.Logger
	Sch      *types.Scheduler
	Jobs     []*types.Job
	cron     gocron.Scheduler
	tasks    map[string]gocron.Task
	signalCh chan os.Signal
	exitCh   chan int
}

func Create(db storage.Provider, log *logr.Logger) Scheduler {
	return &SchedulerProcess{Db: db, Log: log}
}

func (s *SchedulerProcess) Run(label string) {
	var err error
	s.Label = label
	s.InitDb()
	s.CreateInstance()
	s.signalCh = make(chan os.Signal, 2)
	go s.HandleSignals()
	s.cron, err = gocron.NewScheduler()
	if err != nil {
		s.Log.Error(err)
		return
	}
	s.AddSchedulerJobs()
	s.cron.Start()
	<-s.exitCh
	err = s.cron.Shutdown()
	if err != nil {
		s.Log.Error(err)
		return
	}
	s.Log.Info("Scheduler ", s.Sch.ID, " successfully stopped.")
}

func (s *SchedulerProcess) InitDb() {
	s.Db.Init()
	s.Db.Seed()
	s.Db.Free()
}

func (s *SchedulerProcess) CreateInstance() {
	s.Db.Init()
	defer s.Db.Free()
	sch := &types.Scheduler{
		Label: s.Label,
	}
	sch, err := s.Db.InsertScheduler(sch)
	if err != nil {
		s.Log.Error(err)
		return
	}
	s.Sch = sch
	s.Log.Info("Scheduler ", sch.ID, " with label ", s.Label, " started.")
}

func (s *SchedulerProcess) AddSchedulerJobs() {
	_, err := s.cron.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(ping_job, s),
	)
	if err != nil {
		s.Log.Error(err.Error())
	}
	_, err = s.cron.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(get_job, s),
	)
	if err != nil {
		s.Log.Error(err.Error())
	}
}

func (s *SchedulerProcess) HandleSignals() {
	s.Log.Debug("Listening for signals...")
	sig := <-s.signalCh
	switch sig {
	case os.Interrupt, syscall.SIGTERM:
		s.exitCh <- 1
	}
}

func (s *SchedulerProcess) RefreshJobs() {
	s.Db.Init()
	defer s.Db.Free()
	for _, j := range s.Jobs {
		_, ok := s.tasks[j.ID]
		if !ok && len(j.SchedulerId) == 0 {
			s.StartTask(j)
		}
	}
}

func (s *SchedulerProcess) StartTask(job *types.Job) {
	t, err := s.Db.GetTargetByID(job.TargetId, false)
	if err != nil {
		s.Log.Error("[StartTask] Error creating task for target ", job.TargetId, ". ", err.Error())
		return
	}
	m, err := s.Db.GetMonitorByID(t.Monitor.MonitorID)
	if err != nil {
		s.Log.Error("[StartTask] Error getting monitor ", t.Monitor.MonitorID, ". ", err.Error())
		return
	}
	var runtime RuntimeJob
	switch {
	case m.Definition.Shell != nil:
		runtime = shell_job
	case m.Definition.Docker != nil:
		runtime = docker_job
	default:
		s.Log.Error("[StartTask] Invalid monitor.")
		return
	}
	_, err = s.cron.NewJob(
		gocron.CronJob(m.Schedule, false),
		gocron.NewTask(runtime, s, job, m, t.Monitor.Params),
	)
	if err != nil {
		s.Log.Error("[StartTask] Could not start job ", job.ID, ". ", err.Error())
	}
}
