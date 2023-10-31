package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/current", currentVersionHandler)

	http.HandleFunc("/update", updateHandler)
	
	http.HandleFunc("/update/", updateHandler)

	http.HandleFunc("/listversions", listVersionsHandler)

	// Start the HTTP server on port 8080
	fmt.Println("Server listening on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
