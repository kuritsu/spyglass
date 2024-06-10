package scheduler

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kuritsu/spyglass/api/types"
)

func shell_job(s *SchedulerProcess, job *types.Job, monitor *types.Monitor, params map[string]string) {
	s.Log.Debug("[shell_job] Starting shell job ", job.ID, " for target ", job.TargetId, "...")
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		s.Log.Error("[shell_job] Error ", job.ID)
	}
	file.WriteString(monitor.Definition.Shell.Command)
	file.Close()
	defer os.Remove(file.Name())
	cmd := exec.Command(monitor.Definition.Shell.Executable, file.Name())
	cmd.Env = append(cmd.Environ(), fmt.Sprintf("TARGET_ID=%v", job.TargetId))
	for k, v := range monitor.Definition.Shell.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", k, v))
	}
	for k, v := range params {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", k, v))
	}
	s.Log.Info("[shell_job] === Output start. Job ", job.ID, " ===")
	stdout, err := cmd.Output()
	if err != nil {
		s.Log.Error(err.Error())
	}
	fmt.Println(string(stdout))
	s.Log.Info("[shell_job] === Output end. Job ", job.ID, " ===")
	s.Log.Debug("[shell_job] Stopped shell job ", job.ID)
}
