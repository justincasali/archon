package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Task is the basic config unit
type Task struct {
	Resource   string            `yaml:"resource"`
	Action     string            `yaml:"action"`
	Parameters map[string]string `yaml:"parameters"`
	Payload    string            `yaml:"payload"`
}

// Run builds and sends the task to the corresponding server agent
func (task Task) Run(server, token string) error {

	// build request url
	requestURL := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(server, port),
		Path:   task.Resource,
	}

	// build request query
	query := make(url.Values)
	for key, value := range task.Parameters {
		query.Set(key, value)
	}
	requestURL.RawQuery = query.Encode()

	// body for request
	var requestBody io.Reader

	// check if payload exists in task
	if task.Payload != "" {

		// open payload file
		file, err := os.Open(task.Payload)
		if err != nil {
			return err
		}
		defer file.Close()

		// assign file contents to body
		requestBody = file

	}

	// build request for task
	request, err := http.NewRequest(strings.ToUpper(task.Action), requestURL.String(), requestBody)
	if err != nil {
		return err
	}

	// add authentication token to request header
	request.Header.Set("token", token)

	// send request with default http client
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// throw error if status code is not OK
	if response.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(response.StatusCode))
	}

	// success
	return nil

}

// String formats task compactly for display
func (task Task) String() string {
	return fmt.Sprintf("%s.%s", task.Resource, task.Action)
}
