package main

import (
	"context"
	"fmt"
	container "github.com/docker/docker/api/types/container"
	"log"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func runContainer(cli *client.Client, imagename, containername string, inputEnv []string) error {
	port := "8000"   // Container's port
	hostPort := "80" // Host's port

	newport, err := nat.NewPort("tcp", port)
	if err != nil {
		return fmt.Errorf("unable to create docker port: %v", err)
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			newport: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			},
		},
		RestartPolicy: container.RestartPolicy{
			Name: "always",
		},
		LogConfig: container.LogConfig{
			Type:   "json-file",
			Config: map[string]string{},
		},
	}

	config := &container.Config{
		Image:        imagename,
		Env:          inputEnv,
		ExposedPorts: nat.PortSet{newport: struct{}{}},
		Hostname:     fmt.Sprintf("%s-hostnameexample", imagename),
	}

	cont, err := cli.ContainerCreate(
		context.Background(),
		config,
		hostConfig,
		nil, // No need for networking config in this case
		nil, // No need for platform spec
		containername,
	)
	if err != nil {
		return fmt.Errorf("unable to create container: %v", err)
	}

	err = cli.ContainerStart(context.Background(), cont.ID, container.StartOptions{})
	if err != nil {
		return fmt.Errorf("unable to start container: %v", err)
	}

	log.Printf("Container %s is created and running", cont.ID)
	return nil
}

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Unable to create docker client: %v", err)
	}

	imagename := "this_is_an_image_name"
	containername := "this_is_an_image_name"
	inputEnv := []string{"LISTENINGPORT=8000"} // Set environment variable
	err = runContainer(cli, imagename, containername, inputEnv)
	if err != nil {
		log.Fatalf("Error running container: %v", err)
	}
}