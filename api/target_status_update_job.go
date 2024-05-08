package api

import (
	logr "github.com/sirupsen/logrus"

	"github.com/kuritsu/spyglass/api/storage"
	"github.com/kuritsu/spyglass/api/types"
)

type StatusUpdateJob struct {
	db         storage.Provider
	log        *logr.Logger
	StatusChan chan string
	ExitChan   chan int
}

func NewStatusUpdateJob(log *logr.Logger) *StatusUpdateJob {
	// TODO: Improve with dependency injection of db field
	return &StatusUpdateJob{storage.CreateProviderFromConf(log), log, make(chan string), make(chan int)}
}

func (j *StatusUpdateJob) Run() {
	j.log.Debug("[StatusUpdateJob] Job started.")
	for {
		id := <-j.StatusChan
		j.log.Debug("[StatusUpdateJob] Received update request ", id)
		err := j.updateTargetStatus(id)
		if err != nil {
			j.log.Error("[StatusUpdateJob] ", err)
		}
	}
}

func (j *StatusUpdateJob) updateTargetStatus(id string) error {
	j.db.Init()
	defer j.db.Free()
	j.log.Debug("[StatusUpdateJob] Getting target ", id)
	target, err := j.db.GetTargetByID(id, true)
	if err != nil {
		return err
	}
	if len(target.Children) == 0 {
		j.log.Info("[StatusUpdateJob] No children for ", id)
		return nil
	}
	status := 0
	for _, c := range target.Children {
		status += c.Status
	}
	j.log.Debug("Status", status)
	status = int(100.0 * float64(status) / float64(len(target.Children)*100))
	j.log.Debug("[StatusUpdateJob] Updating target ", id, " with status ", status)
	_, err = j.db.UpdateTargetStatus(target, &types.TargetPatch{
		Status:            status,
		StatusDescription: "Updated automatically",
	})
	return err
}
