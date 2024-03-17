package api

import (
	"ContainerDeployerManager/internal/deployerManager"
	deployer "ContainerDeployerManager/proto"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

// DeployRequest represents the request payload for deploying a container
type DeployRequest struct {
	ImageName   string            `json:"imageName"`
	Args        []string          `json:"args"`
	EnvVars     map[string]string `json:"envVars"`
	ProcessPort string            `json:"processPort"`
}

// DeployResponse represents the response payload after attempting to deploy a container
type DeployResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// DeployContainerHandler handles the HTTP request to deploy a container
func DeployContainerHandler(w http.ResponseWriter, r *http.Request) {
	var req DeployRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.ProcessPort == "" {
		log.Println("No process Port identified")
		dockerClient, err := deployerManager.NewDockerClient()
		if err != nil {
			http.Error(w, "Failed to create Docker client: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = dockerClient.DeployContainer(r.Context(), req.ImageName, req.Args, req.EnvVars)
		if err != nil {
			http.Error(w, "Failed to deploy container: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		callGRPCProcess(req)
	}

	response := DeployResponse{
		Status:  "Success",
		Message: "Container deployed successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func callGRPCProcess(req DeployRequest) {
	// Server address
	serverAddr := "localhost:" + req.ProcessPort

	// Set up a connection to the server.
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	defer conn.Close()
	c := deployer.NewDeployerServiceClient(conn)

	// Prepare and send the DeployContainer request
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.DeployContainer(ctx, &deployer.DeployRequest{
		ImageName: req.ImageName,
		Args:      req.Args,
		EnvVars:   req.EnvVars,
	})

	if err != nil {
		log.Printf("could not deploy: %v", err)
	}
	log.Printf("Deployment response by Grpc: %s", r.Message)
}

// SetupRoutes sets up the router and routes for the HTTP server
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/deploy", DeployContainerHandler).Methods("POST")
	return r
}
