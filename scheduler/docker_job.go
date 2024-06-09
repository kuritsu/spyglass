package scheduler

import (
	"github.com/kuritsu/spyglass/api/types"
)

func docker_job(s *SchedulerProcess, job *types.Job, monitor *types.Monitor, params map[string]string) {
	// db := storage.CreateProviderFromConf(s.Log)
	s.Log.Debug("Starting docker job ", job.ID, " for target ", job.TargetId, "...")
	//TODO: Run docker
	s.Log.Debug("Stopped docker job ", job.ID)
}
