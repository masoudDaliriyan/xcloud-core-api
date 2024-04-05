package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func buildImage(cli *client.Client, tags []string, dockerfile string) error {
	ctx := context.Background()

	// Open the Dockerfile
	dockerFileReader, err := os.Open(dockerfile)
	if err != nil {
		return err
	}
	defer dockerFileReader.Close()

	// Read the actual Dockerfile content
	dockerFileContent, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		return err
	}

	// Create a buffer
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Make a TAR header for the Dockerfile
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(dockerFileContent)),
	}
	if err := tw.WriteHeader(tarHeader); err != nil {
		return err
	}

	// Write the Dockerfile content to the TAR file
	if _, err := tw.Write(dockerFileContent); err != nil {
		return err
	}

	// Close the TAR writer
	if err := tw.Close(); err != nil {
		return err
	}

	// Create a TAR reader from the buffer
	tarReader := bytes.NewReader(buf.Bytes())

	// Define the build options
	buildOptions := types.ImageBuildOptions{
		Tags:       tags,
		Dockerfile: dockerfile,
		Context:    tarReader,
		Remove:     true,
	}

	// Build the image
	imageBuildResponse, err := cli.ImageBuild(ctx, tarReader, buildOptions)
	if err != nil {
		return err
	}
	defer imageBuildResponse.Body.Close()

	// Print the build output
	fmt.Println("Building Docker image...")
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Unable to create docker client: %s", err)
	}

	// Define image tags and Dockerfile location
	tags := []string{"this_is_an_image_name"}
	dockerfile := "./repos/simple-http-server/Dockerfile"

	// Build the Docker image
	err = buildImage(cli, tags, dockerfile)
	if err != nil {
		log.Fatalf("Error building Docker image: %s", err)
	}
}