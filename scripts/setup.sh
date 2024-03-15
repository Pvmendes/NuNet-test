#!/bin/bash

# Exit on any error
set -e

# Check for Docker installation
if ! [ -x "$(command -v docker)" ]; then
  echo "Error: Docker is not installed." >&2
  exit 1
fi

# Ensure the Docker daemon is running
if ! docker info >/dev/null 2>&1; then
    echo "Error: Docker daemon is not running." >&2
    exit 1
fi

# # Install Go dependencies
# echo "Installing Go dependencies..."
# go mod tidy

# # Building the server
# echo "Building the server binary..."
# go build -o server ./cmd/server

# # Building the client
# echo "Building the client binary..."
# go build -o client ./cmd/client

# echo "Setup completed successfully."

# Building Docker image
echo "Building Docker image for the application..."
docker build -t container-deployer-manager .

# Optionally, run the Docker container immediately after build
echo "Running the Docker container..."
docker run -d -p 8080:8080 container-deployer-manager

echo "Setup completed successfully. Application is now running in a Docker container."
