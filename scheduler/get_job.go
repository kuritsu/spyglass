package scheduler

import "github.com/kuritsu/spyglass/api/storage"

func get_job(s *SchedulerProcess) {
	db := storage.CreateProviderFromConf(s.Log)
	db.Init()
	defer db.Free()
	s.Log.Debug("[get_job] Getting list of jobs for label ", s.Label, "...")
	jobs, err := db.GetAllJobsFor(s.Label)
	if err != nil {
		s.Log.Error("[get_job] ", err.Error())
		return
	}
	s.Jobs = jobs
	s.Log.Debug("[get_job] Refreshing ", len(jobs), " jobs...")
	s.RefreshJobs()
}
