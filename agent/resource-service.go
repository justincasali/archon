package main

import (
	"net/http"
	"os/exec"
)

// StartService starts a linux service on the server
func StartService(w http.ResponseWriter, r *http.Request) {

	// get service from query
	svc := r.URL.Query().Get("service")
	if svc == "" {
		http.Error(w, "invalid service parameter", http.StatusBadRequest)
		return
	}

	// start service
	startOutput, err := exec.CommandContext(r.Context(), "systemctl", "start", svc).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(startOutput)
		return
	}

	// write output
	w.Write(startOutput)

}

// StopService stops a linux service on the server
func StopService(w http.ResponseWriter, r *http.Request) {

	// get service from query
	svc := r.URL.Query().Get("service")
	if svc == "" {
		http.Error(w, "invalid service parameter", http.StatusBadRequest)
		return
	}

	// stop service
	stopOutput, err := exec.CommandContext(r.Context(), "systemctl", "stop", svc).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(stopOutput)
		return
	}

	// write output
	w.Write(stopOutput)

}
