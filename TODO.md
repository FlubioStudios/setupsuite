# DOCKER LOG FIX
  create file /etc/docker/daemon.json
  ```
  {
    "log-driver": "local",
    "log-opts": {
    "max-size": "20m",
    "max-file": "5"
    }
  }
  ```
  restart docker service
