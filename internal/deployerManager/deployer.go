package deployerManager

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerClient wraps the Docker API client
type DockerClient struct {
	cli *client.Client
}

// NewDockerClient creates a new DockerClient
func NewDockerClient() (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerClient{cli: cli}, nil
}

// DeployContainer deploys a container based on the specified options
func (dc *DockerClient) DeployContainer(ctx context.Context, image string, args []string, envVars map[string]string) error {
	fmt.Println("Start Deploy Container")
	// Pull the image
	out, err := dc.cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Pull the image")
	io.Copy(os.Stdout, out) // Optionally, output the pull progress to stdout

	// Convert envVars map to slice of strings
	env := make([]string, 0, len(envVars))
	for k, v := range envVars {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	fmt.Println("Convert envVars map to slice of strings")

	// Create the container
	resp, err := dc.cli.ContainerCreate(ctx, &container.Config{
		Image: image,
		Cmd:   args,
		Env:   env,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}
	fmt.Println("Create the container")

	// Start the container
	if err := dc.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	fmt.Printf("Container %s is deployed\n", resp.ID)
	return nil
}
