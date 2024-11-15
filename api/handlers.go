package api

import (
	"encoding/json"
	"net/http"
	"retail_pulse/internal/service"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetJobInfoHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	jobID := query.Get("jobid")

	if jobID == "" {
		http.Error(w, "jobid is a required field", 400)
		return
	}

	id, err := primitive.ObjectIDFromHex(jobID)

	if err != nil {
		http.Error(w, "invalid jobid", 400)
		return
	}

	svs := service.NewStoresVisitService()

	status, errMssg, failedStoreID, err := svs.GetStatusAndErrorByID(id)

	if err != nil {
		http.Error(w, "jobid does not exist", http.StatusBadRequest)
	} else {
		if status == "completed" || status == "ongoing" {
			res := struct {
				Status string `json:"status"`
				JobID  string `json:"job_id"`
			}{
				Status: status,
				JobID:  jobID,
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		} else {
			type ErrStruct struct {
				StoreID string `json:"store_id"`
				Error   string `json:"error"`
			}

			res := struct {
				Status string    `json:"status"`
				JobID  string    `json:"job_id"`
				Error  ErrStruct `json:"error"`
			}{
				Status: status,
				JobID:  jobID,
				Error:  ErrStruct{StoreID: failedStoreID, Error: errMssg},
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res)
		}
	}
}
