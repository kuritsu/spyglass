package scheduler

import "github.com/kuritsu/spyglass/api/storage"

func ping_job(s *SchedulerProcess) {
	db := storage.CreateProviderFromConf(s.Log)
	db.Init()
	defer db.Free()
	sch, err := db.UpdateScheduler(s.Sch)
	if err != nil {
		s.Log.Error("[ping_job] ", err.Error())
		return
	}
	s.Sch = sch
	s.Log.Debug("[ping_job] Scheduler ", s.Sch.ID, " updated at ", s.Sch.LastPing)
}
