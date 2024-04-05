package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"log"
)

// Stop and remove a container
func stopAndRemoveContainer(client *client.Client, containername string) error {
	ctx := context.Background()

	if err := client.ContainerStop(ctx, containername, container.StopOptions{}); err != nil {
		log.Printf("Unable to stop container %s: %s", containername, err)
	}

	removeOptions := container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := client.ContainerRemove(ctx, containername, removeOptions); err != nil {
		log.Printf("Unable to remove container: %s", err)
		return err
	}

	return nil
}

func main() {
	client, err := client.NewEnvClient()
	if err != nil {
		fmt.Printf("Unable to create docker client: %s", err)
	}

	// Stops and removes a container
	stopAndRemoveContainer(client, "this_is_an_image_name")
}