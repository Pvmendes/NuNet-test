package config

import (
	"ContainerDeployerManager/internal/api"
	"ContainerDeployerManager/internal/deployerManager"
	deployer "ContainerDeployerManager/proto"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"

	"google.golang.org/grpc"
)

// AppConfig holds the configuration settings for the application.
type AppConfig struct {
	Port           int    `json:"serverPort"`
	GrpcPort       string `json:"serverGrpcPort"`
	DockerEndpoint string `json:"dockerEndpoint"`
}

type AppConfigBoth struct {
	ServerPort     int    `json:"serverPort"`
	ServerGrpcPort string `json:"serverGrpcPort"`
	ClientPort     int    `json:"clientPort"`
	ClientGrpcPort string `json:"clientGrpcPort"`
	DockerEndpoint string `json:"dockerEndpoint"`
}

// LoadConfig loads configuration settings from environment variables.
func LoadAppConfigBothConfig() *AppConfigBoth {
	config := &AppConfigBoth{
		ServerPort:     getEnvAsInt("SERVER_PORT", 8080),
		ServerGrpcPort: getEnvAsString("SERVER_GRPC_PORT", ":50050"),
		ClientPort:     getEnvAsInt("CLIENT_PORT", 8082),
		ClientGrpcPort: getEnvAsString("CLIENT_GRPC_PORT", ":50052"),
		DockerEndpoint: getEnvAsString("DOCKER_ENDPOINT", "unix:///var/run/docker.sock"),
	}
	return config
}

// Helper function to read an environment variable as a string or return a default value
func getEnvAsString(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as an integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err == nil {
			return intValue
		}
		log.Printf("Warning: Could not parse %s as integer, defaulting to %d", key, defaultValue)
	}
	return defaultValue
}

func CreateAppConfig(port int, grpcPort, dockerEndpoint string) *AppConfig {
	return &AppConfig{
		Port:           port,
		GrpcPort:       grpcPort,
		DockerEndpoint: dockerEndpoint,
	}
}

type client struct {
	deployer.UnimplementedDeployerServiceServer
}

func (s *client) DeployContainer(ctx context.Context, req *deployer.DeployRequest) (*deployer.DeployResponse, error) {
	dockerClient, err := deployerManager.NewDockerClient()
	if err != nil {
		return nil, err
	}

	err = dockerClient.DeployContainer(ctx, req.ImageName, req.Args, req.EnvVars)
	if err != nil {
		return &deployer.DeployResponse{Status: "Failure", Message: err.Error()}, nil
	}

	return &deployer.DeployResponse{Status: "Success", Message: "Container deployed successfully"}, nil
}

func StartLoadServers(cfg *AppConfig) {
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", cfg.GrpcPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		server := grpc.NewServer()
		deployer.RegisterDeployerServiceServer(server, &client{})

		log.Printf("Starting GRPC server listening at %v", lis.Addr())

		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start the HTTP server on a different port
	go func() {
		defer wg.Done()
		// Initialize and start the HTTP server
		r := api.SetupRoutes()

		log.Printf("Starting HTTP server on port %d", cfg.Port)

		if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	wg.Wait()
}
