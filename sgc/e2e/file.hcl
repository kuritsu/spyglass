# This is a comment

monitor "monitors.mymonitor2" {
  type = "docker"
  schedule = "* * * * *"
  definition {
    docker {
      image = "nginx:latest"
      docker_env = {
        var1 = "var1"
        var2 = "var2"
      }
    }
  }
  readers = ["ex"]
}

target "target" {
  description = "this is my target"
	url = "https://mytarget.url"
	view {
    image_big = "http://mytarget.url/big"
  }
	status = 6
	status_description = "Progress of task"
	critical = true
	monitor {
    monitor_id = "monitors.mymonitor2"
  }
	writers = ["user@email.com"]
}
