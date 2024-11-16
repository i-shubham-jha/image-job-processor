package main

import (
	"flag"
	"fmt"
	"net/http"
	"retail_pulse/api"
	"retail_pulse/internal/store"

	"github.com/gorilla/mux"
)

func main() {
	// reading cmd args
	port := flag.Int("p", 8080, "Port number to start server on")
	file := flag.String("f", "StoreMasterAssignment.csv", "File name to read")

	// parse the command line flags
	flag.Parse()

	// set csv file
	store.CsvFilePath = *file

	// routes and handlers
	r := mux.NewRouter()
	r.HandleFunc("/api/status", api.GetJobInfoHandler)
	r.HandleFunc("/api/submit", api.SubmitJobHandler).Methods("POST")

	// start server
	http.ListenAndServe(fmt.Sprintf("localhost:%d", *port), r)
}
