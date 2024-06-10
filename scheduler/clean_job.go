package scheduler

import "github.com/kuritsu/spyglass/api/storage"

func clean_job(s *SchedulerProcess) {
	db := storage.CreateProviderFromConf(s.Log)
	db.Init()
	defer db.Free()
	s.Log.Debug("[clean_job] Getting all inactive schedulers...")
	schedulers, err := db.GetAllInactiveSchedulers()
	if err != nil {
		s.Log.Error("[clean_job] ", err.Error())
		return
	}
	s.Log.Debug("[clean_job] Cleaning ", len(schedulers), " schedulers...")
	for _, sch := range schedulers {
		count, err := db.UnassignJobs(sch.ID)
		if err != nil {
			s.Log.Error("[clean_job] Error unassigning jobs. ", err.Error())
			return
		}
		s.Log.Debug("[clean_job] Updated ", count, " jobs with scheduler ", sch.ID, ".")
		err = db.DeleteScheduler(sch.ID)
		if err != nil {
			s.Log.Error("[clean_job] Error deleting scheduler. ", err.Error())
			return
		}
	}
}
