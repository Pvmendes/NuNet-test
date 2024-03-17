# Container Deployer/Manager

Container Deployer/Manager is a tool designed to facilitate the deployment and management of containers across networked machines. It allows users to deploy a container on a remote machine through a simple API or user interface.

## Features

- Deploy containers on remote machines via a RESTful API or gRPC.
- Support for Docker containers.
- Simple CLI for managing deployments.
- Secure communication between client and server.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

What things you need to install the software and how to install them:

- Docker
- Go (version 1.18 or later)

### Running the Setup Script

Open a terminal in the project's root directory and execute the following command:

```bash
./setup.sh
```

### Deploy a Container

To deploy a container, send a POST request to the `/deploy` endpoint with the image name, arguments, and environment variables for the container. 
Replace `localhost:8080` with your server's address if different.

```bash
curl -X POST http://localhost:8080/deploy \
-H "Content-Type: application/json" \
-d '{
  "imageName": "busybox",
  "args": ["echo", "Hello World"],
  "envVars": {
    "EXAMPLE_VAR": "value"
  }
}'
```
This request will deploy a container based on the BusyBox image that echoes "Hello World". 
The response from the server will include the status of the deployment and any messages related to the deployment process.

If the deployment is intended to be executed by the client, please include "processPort": "50052" in the corresponding JSON file.

```bash
curl -X POST http://localhost:8080/deploy \
-H "Content-Type: application/json" \
-d '{
  "imageName": "busybox",
  "args": ["echo", "Hello World"],
  "envVars": {
    "EXAMPLE_VAR": "value"
  },
  "processPort":"50052"
}'
```

#### Expected Response
Upon successful deployment, you should receive a response similar to the following:
```bash
{
  "status": "Success",
  "message": "Container deployed successfully"
}
```