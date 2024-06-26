package scheduler

import (
	"os"
	"os/signal"
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
	tasks    map[string]gocron.Job
	signalCh chan os.Signal
	exitCh   chan int
}

func Create(db storage.Provider, log *logr.Logger) Scheduler {
	return &SchedulerProcess{Db: db, Log: log, tasks: make(map[string]gocron.Job)}
}

func (s *SchedulerProcess) Run(label string) {
	var err error
	s.Label = label
	s.InitDb()
	s.CreateInstance()
	s.signalCh = make(chan os.Signal, 1)
	signal.Notify(s.signalCh, os.Interrupt, syscall.SIGTERM)
	go s.HandleSignals()
	s.cron, err = gocron.NewScheduler()
	if err != nil {
		s.Log.Error(err)
		return
	}
	s.AddSchedulerJobs()
	s.cron.Start()
	<-s.exitCh
}

func (s *SchedulerProcess) InitDb() {
	s.Db.Init()
	s.Db.Seed()
	s.Db.Free()
}

func (s *SchedulerProcess) Free() {
	err := s.cron.Shutdown()
	if err != nil {
		s.Log.Error(err)
		return
	}
	s.Log.Info("Scheduler ", s.Sch.ID, " successfully stopped.")
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
	s.addJob(ping_job, 1*time.Minute)
	s.addJob(get_job, 1*time.Minute)
	s.addJob(clean_job, 5*time.Minute)
}

func (s *SchedulerProcess) addJob(function any, duration time.Duration) {
	_, err := s.cron.NewJob(
		gocron.DurationJob(duration),
		gocron.NewTask(function, s),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)
	if err != nil {
		s.Log.Error(err.Error())
		os.Exit(1)
		return
	}
}

func (s *SchedulerProcess) HandleSignals() {
	s.Log.Debug("Listening for signals...")
	<-s.signalCh
	s.Free()
	os.Exit(1)
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
	s.Log.Debug("[StartTask] Starting job ", job.ID, " with cron ", m.Schedule)
	cronjob, err := s.cron.NewJob(
		gocron.CronJob(m.Schedule, false),
		gocron.NewTask(runtime, s, job, m, t.Monitor.Params),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)
	if err != nil {
		s.Log.Error("[StartTask] Could not start job ", job.ID, ". ", err.Error())
		return
	}
	s.tasks[job.ID] = cronjob
	job.SchedulerId = s.Sch.ID
	job, err = s.Db.UpdateJob(job)
	if err != nil {
		s.Log.Error("[StartTask] Could not update job ", job.ID, ". ", err.Error())
		return
	}
	s.Log.Debug("Job ", job.ID, " updated.")
}
