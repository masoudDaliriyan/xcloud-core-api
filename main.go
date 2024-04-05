package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {
	// Specify the directory containing your Dockerfile
	buildContextDir := "./repos/simple-http-server"
	dockerfilePath := "Dockerfile" // Assuming the Dockerfile is named "Dockerfile" in the build context directory

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		return
	}

	// Open the project directory
	buildContext, err := os.Open(buildContextDir)
	if err != nil {
		fmt.Printf("Error opening build context: %v\n", err)
		return
	}
	defer buildContext.Close()

	// Build Docker image
	buildOptions := types.ImageBuildOptions{
		Tags:       []string{"simple-http-server:latest"},
		Dockerfile: dockerfilePath,
		Context:    buildContext,
	}
	resp, err := cli.ImageBuild(context.Background(), buildContext, buildOptions)
	if err != nil {
		fmt.Printf("Error building Docker image: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Print build output
	fmt.Println("Building Docker image...")
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Printf("Error printing build output: %v\n", err)
		return
	}

	fmt.Println("Docker image built successfully")
}