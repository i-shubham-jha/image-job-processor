package main

import (
	"flag"
	"fmt"
	"image-job-processor/api"
	"image-job-processor/internal/logger"
	"image-job-processor/internal/service"
	"image-job-processor/internal/store"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	logger := logger.GetLogger()

	// reading cmd args
	port := flag.Int("p", 8080, "Port number to start server on")
	file := flag.String("f", "StoreMasterAssignment.csv", "File name to read")

	// parse the command line flags
	flag.Parse()

	// set csv file
	store.CsvFilePath = *file
	logger.Log(fmt.Sprintf("Reading csv file %s", *file))
	store.NewStoreManager()

	// establish connection to mongodb
	service.NewStoresVisitService()

	// routes and handlers
	r := mux.NewRouter()
	r.HandleFunc("/api/status", api.GetJobInfoHandler).Methods("GET")
	r.HandleFunc("/api/submit", api.SubmitJobHandler).Methods("POST")

	// start server
	logger.Log(fmt.Sprintf("Starting server on port %v", *port))
	logger.Log("-------------INIT DONE-------------")
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", *port), r)
}
