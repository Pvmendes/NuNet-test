package main

import (
	"ContainerDeployerManager/internal/config"
)

func main() {
	cfg := config.LoadAppConfigBothConfig()

	config.StartLoadServers(config.CreateAppConfig(cfg.ClientPort, cfg.ClientGrpcPort, cfg.DockerEndpoint))
}
