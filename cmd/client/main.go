package main

import (
	deployer "ContainerDeployerManager/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// Set up a connection to the server.
	address := "localhost:50051" // Change this to your server's address
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := deployer.NewDeployerServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	imageName := "busybox" // Example image name
	args := []string{"echo", "Hello from the deployed container!"}
	envVars := map[string]string{"EXAMPLE_VAR": "VALUE"}

	r, err := c.DeployContainer(ctx, &deployer.DeployRequest{
		ImageName: imageName,
		Args:      args,
		EnvVars:   envVars,
	})
    
	if err != nil {
		log.Fatalf("could not deploy: %v", err)
	}
	log.Printf("Deployment status: %s, message: %s", r.Status, r.Message)
}
