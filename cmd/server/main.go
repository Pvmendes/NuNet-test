package main

import (
	"ContainerDeployerManager/internal/deployerManager"
	deployer "ContainerDeployerManager/proto"
	"context"
	"log"
	"net"

	"ContainerDeployerManager/internal/api"
	"ContainerDeployerManager/internal/config"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
)

type server struct {
	deployer.UnimplementedDeployerServiceServer
}

func (s *server) DeployContainer(ctx context.Context, req *deployer.DeployRequest) (*deployer.DeployResponse, error) {
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

func main() {
	cfg := config.LoadConfig()

	// Initialize and start the HTTP server
	r := api.SetupRoutes()
	log.Printf("Starting HTTP server on port %d", cfg.ServerPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ServerPort), r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	deployer.RegisterDeployerServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
