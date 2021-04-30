package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"path"
	"strings"
)

// CreateFile creates a new file on the server from the request's body
func CreateFile(w http.ResponseWriter, r *http.Request) {

	// get file from query; ensure path is absolute
	file := r.URL.Query().Get("file")
	if !path.IsAbs(file) {
		http.Error(w, "invalid file parameter", http.StatusBadRequest)
		return
	}

	// create command
	create := exec.CommandContext(r.Context(), "tee", file)

	// grab command input
	input, err := create.StdinPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// stream request body into command in background
	go func() {
		defer input.Close()
		io.Copy(input, r.Body)
	}()

	// create file
	createOutput, err := create.CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(createOutput)
		return
	}

	// check if mode in query
	if mode := r.URL.Query().Get("mode"); mode != "" {

		// chmod file
		chmodOutput, err := exec.CommandContext(r.Context(), "chmod", mode, file).CombinedOutput()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(chmodOutput)
			return
		}

	}

	// check if owner or group in query
	if owner, group := r.URL.Query().Get("owner"), r.URL.Query().Get("group"); owner != "" || group != "" {

		// join group with owner if exits
		if group != "" {
			owner = strings.Join([]string{owner, group}, ":")
		}

		// chown file
		chownOutput, err := exec.CommandContext(r.Context(), "chown", owner, file).CombinedOutput()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(chownOutput)
			return
		}

	}

	// write output as md5 hash of the file created
	fmt.Fprintf(w, "%x\n", md5.Sum(createOutput))

}

// DestroyFile removes a file on the server
func DestroyFile(w http.ResponseWriter, r *http.Request) {

	// get file from query; ensure path is absolute
	file := r.URL.Query().Get("file")
	if !path.IsAbs(file) {
		http.Error(w, "invalid file parameter", http.StatusBadRequest)
		return
	}

	// remove file
	removeOutput, err := exec.CommandContext(r.Context(), "rm", "-f", file).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(removeOutput)
		return
	}

	// write output
	w.Write(removeOutput)

}
