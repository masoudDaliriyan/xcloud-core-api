package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

func runContainer(imagename, containername string, inputEnv []string) error {
    // Define the Docker Engine API endpoint
    url := "http://127.0.0.1:12345/containers/create"

    // Prepare the request payload
    body := map[string]interface{}{
        "Image":        imagename,
        "Env":          inputEnv,
        "ExposedPorts": map[string]struct{}{"8000/tcp": {}},
        "HostConfig": map[string]interface{}{
            "PortBindings": map[string][]map[string]string{
                "8000/tcp": {
                    {"HostIP": "0.0.0.0", "HostPort": "80"},
                },
            },
            "RestartPolicy": map[string]interface{}{"Name": "always"},
            "LogConfig":     map[string]string{"Type": "json-file"},
        },
    }

    // Convert the payload to JSON
    payload, err := json.Marshal(body)
    if err != nil {
        return fmt.Errorf("unable to marshal JSON payload: %v", err)
    }

    // Send the HTTP POST request to create the container
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return fmt.Errorf("unable to create container: %v", err)
    }
    defer resp.Body.Close()

    // Read the response body
    respBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("unable to read response body: %v", err)
    }

    // Check if the request was successful
    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("failed to create container: %s", respBody)
    }

    // Extract the container ID from the response
    var createResp map[string]interface{}
    if err := json.Unmarshal(respBody, &createResp); err != nil {
        return fmt.Errorf("unable to unmarshal response body: %v", err)
    }

    containerID, ok := createResp["Id"].(string)
    if !ok {
        return fmt.Errorf("unable to extract container ID from response")
    }

    log.Printf("Container %s is created and running", containerID)
    return nil
}

func main() {
    imagename := "this_is_an_image_name"
    containername := "this_is_an_image_name"
    inputEnv := []string{"LISTENINGPORT=8000"} // Set environment variable

    err := runContainer(imagename, containername, inputEnv)
    if err != nil {
        log.Fatalf("Error running container: %v", err)
    }
}
