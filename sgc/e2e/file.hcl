# This is a comment

monitor "mymonitors.mymonitor-2" {
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
