package api

import (
	"encoding/json"
	"fmt"
	"image-job-processor/internal/job"
	"image-job-processor/internal/model"
	"image-job-processor/internal/service"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetJobInfoHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	jobID := query.Get("jobid")

	if jobID == "" {
		sendErrBack("jobd is a required field", w)
		return
	}

	id, err := primitive.ObjectIDFromHex(jobID)

	if err != nil {
		sendErrBack("invalid jobid", w)
		return
	}

	svs := service.NewStoresVisitService()

	status, errMssg, failedStoreID, err := svs.GetStatusAndErrorByID(id)

	if err != nil {
		sendErrBack("jobid does not exist", w)
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

func SubmitJobHandler(w http.ResponseWriter, r *http.Request) {
	var storesVisit model.StoresVisit

	err := json.NewDecoder(r.Body).Decode(&storesVisit)

	if err != nil {
		sendErrBack("JSON decoding error", w)
		return
	}

	// validate data
	err = validateData(&storesVisit)

	if err != nil {
		sendErrBack(err.Error(), w)
		return
	}

	// insert in db with ongoing status
	storesVisit.Status = "ongoing"

	svs := service.NewStoresVisitService()

	id, err := svs.InsertStoresVisit(storesVisit)

	if err != nil {
		sendErrBack(err.Error(), w)
		return
	}

	// launch go routine for processing
	go job.ProcessJob(id, storesVisit)

	// return job id as json
	res := struct {
		JobID string `json:"job_id"`
	}{JobID: id.Hex()}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func sendErrBack(err string, w http.ResponseWriter) {
	type errRes struct {
		Error string `json:"error"`
	}

	res := errRes{Error: err}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
}

func validateData(sv *model.StoresVisit) error {
	if sv.Count < 0 {
		return fmt.Errorf("count can not be less than zero")
	}

	if len(sv.Visits) == 0 {
		return fmt.Errorf("visits should not be emtpy")
	}

	if len(sv.Visits) != sv.Count {
		return fmt.Errorf("count != len(visits)")
	}

	for _, v := range sv.Visits {
		if v.StoreID == "" {
			return fmt.Errorf("store_id is required")
		}
		if v.VisitTime == "" {
			return fmt.Errorf("visit_time is required")
		}
		if len(v.ImageURLs) == 0 {
			return fmt.Errorf("image_url cannot be empty for store_id: %s", v.StoreID)
		}
	}

	return nil
}
