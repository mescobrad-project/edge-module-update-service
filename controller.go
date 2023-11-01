package main

import (
    "fmt"
    "net/http"
	"encoding/json"
	"strings"
)

var imageName = "mescobrad-edge"


func currentVersionHandler(w http.ResponseWriter, r *http.Request) {
    // Call the checkEdgeModulePresence function from edgepresence package
    result, err := checkEdgeModulePresence(imageName)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    // Write the result to the HTTP response
    fmt.Fprintf(w, result)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
    // Update Edge Module
	version := strings.TrimPrefix(r.URL.Path, "/update/")
	var result string
	var err error
	if version == "/update" {
		result, err = updateEdgeModuleLatest(imageName, imageName)
	} else {
    	result, err = updateEdgeModule(imageName, imageName, version)
	}

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

    // Write the result to the HTTP response
    fmt.Fprintf(w, result)
}

func listVersionsHandler(w http.ResponseWriter, r *http.Request) {
    // Update Edge Module
    result, err := listVersions(imageName)

    if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Convert the array of strings to JSON
	responseJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Write the JSON response to the HTTP response writer
	w.Write(responseJSON)
}

