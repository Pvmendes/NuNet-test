package main

import (
	"ContainerDeployerManager/internal/config"
)

func main() {
	cfg := config.LoadAppConfigBothConfig()

	config.StartLoadServers(config.CreateAppConfig(cfg.ServerPort, cfg.ServerGrpcPort, cfg.DockerEndpoint))
}
