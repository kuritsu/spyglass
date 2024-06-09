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

type SchedulerProcess struct {
	Label    string
	Db       storage.Provider
	Log      *logr.Logger
	Sch      *types.Scheduler
	cron     gocron.Scheduler
	signalCh chan os.Signal
	exitCh   chan int
}

func Create(db storage.Provider, log *logr.Logger) Scheduler {
	return &SchedulerProcess{Db: db, Log: log}
}

func (s *SchedulerProcess) Run(label string) {
	var err error
	s.Label = label
	s.CreateInstance()
	s.signalCh = make(chan os.Signal, 2)
	go s.HandleSignals()
	s.cron, err = gocron.NewScheduler()
	if err != nil {
		s.Log.Error(err)
		return
	}
	s.AddPingJob()
	s.cron.Start()
	<-s.exitCh
	err = s.cron.Shutdown()
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

func (s *SchedulerProcess) AddPingJob() {
	s.cron.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(ping_job, s),
	)
}

func (s *SchedulerProcess) HandleSignals() {
	s.Log.Debug("Listening for signals...")
	sig := <-s.signalCh
	switch sig {
	case os.Interrupt, syscall.SIGTERM:
		s.exitCh <- 1
	}
}
