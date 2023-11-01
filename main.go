package main

import (
	"fmt"
	"net/http"
	"github.com/joho/godotenv"
	"os"
)

var accessKeyID string 
var secretAccessKey string

func init() {
	// Load environment variables from .env file
    if err := godotenv.Load(); err != nil {
        fmt.Println("Error loading .env file")
    }

	accessKeyID = os.Getenv("ACCESSKEYID")
    secretAccessKey = os.Getenv("SECRETACCESSKEY")
}

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