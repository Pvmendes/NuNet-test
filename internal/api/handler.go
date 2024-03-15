package api

import (
	"encoding/json"
	"net/http"
	"ContainerDeployerManager/internal/deployerManager"
	"github.com/gorilla/mux"
)

// DeployRequest represents the request payload for deploying a container
type DeployRequest struct {
	ImageName string            `json:"imageName"`
	Args      []string          `json:"args"`
	EnvVars   map[string]string `json:"envVars"`
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

	response := DeployResponse{
		Status:  "Success",
		Message: "Container deployed successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// SetupRoutes sets up the router and routes for the HTTP server
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/deploy", DeployContainerHandler).Methods("POST")
	return r
}
