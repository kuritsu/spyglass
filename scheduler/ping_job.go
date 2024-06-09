package scheduler

func ping_job(s *SchedulerProcess) {
	s.Db.Init()
	defer s.Db.Free()
	sch, err := s.Db.UpdateScheduler(s.Sch)
	if err != nil {
		s.Log.Error(err.Error())
		return
	}
	s.Sch = sch
	s.Log.Debug("[ping_job] Scheduler ", s.Sch.ID, " updated at ", s.Sch.LastPing)
}
