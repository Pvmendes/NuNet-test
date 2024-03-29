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

# Install Go dependencies
echo "Installing Go dependencies..."
go mod tidy

# Building the server on windows
echo "Building the server binary..."
cd ./cmd/server
GOOS=windows GOARCH=amd64 go build -o ../server.exe

cd ..
cd ..

# Building the client on windows
echo "Building the client binary..."
cd ./cmd/client
GOOS=windows GOARCH=amd64 go build -o ../client.exe

echo "Setup completed successfully."

# # Building Docker image
# echo "Building Docker image for the application..."
# docker build -t container-deployer-manager .

# # Optionally, run the Docker container immediately after build
# echo "Running the Docker container..."
# docker run -d -p 8080:8080 container-deployer-manager

#docker-compose up --build

# echo "Setup completed successfully. Application is now running in a Docker container."
