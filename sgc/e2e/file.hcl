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
  description = "this is my real target"
	url = "https://mytarget.url"
	view {
    image_big = "https://mytarget.url/big"
  }
	status = 55
	status_description = "Progress of task is there"
	critical = true
	monitor {
    monitor_id = "monitors.mymonitor2"
  }
	writers = ["user2@email.com"]
}

# target "target" {
#   description = "this is my target"
# 	url = "https://mytarget2.url"
# 	view {
#     image_small = "http://mytarget.url/small"
#   }
# 	status = 14
# 	status_description = "Progress of task"
# 	critical = true
# 	monitor {
#     monitor_id = "monitors.mymonitor2"
#   }
# 	writers = ["user@email.com"]
# }
