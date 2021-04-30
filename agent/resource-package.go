package main

import (
	"bytes"
	"net/http"
	"os/exec"
)

// InstallPackage adds a debian package to the server
func InstallPackage(w http.ResponseWriter, r *http.Request) {

	// get package from query
	pkg := r.URL.Query().Get("package")
	if pkg == "" {
		http.Error(w, "invalid package parameter", http.StatusBadRequest)
		return
	}

	// update package info
	updateOutput, err := exec.CommandContext(r.Context(), "apt", "update").CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(updateOutput)
		return
	}

	// install package
	installOutput, err := exec.CommandContext(r.Context(), "apt", "install", "-y", pkg).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(installOutput)
		return
	}

	// write output with join
	w.Write(bytes.Join([][]byte{updateOutput, installOutput}, []byte("\n-----\n")))

}

// RemovePackage removes a debian package from the server
func RemovePackage(w http.ResponseWriter, r *http.Request) {

	// get package from query
	pkg := r.URL.Query().Get("package")
	if pkg == "" {
		http.Error(w, "invalid package parameter", http.StatusBadRequest)
		return
	}

	// remove package
	removeOutput, err := exec.CommandContext(r.Context(), "apt", "remove", "-y", pkg).CombinedOutput()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(removeOutput)
		return
	}

	// write output
	w.Write(removeOutput)

}
