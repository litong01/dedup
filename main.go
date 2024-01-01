package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/dedup/utils"
	"github.com/gorilla/mux"
)

func main() {

	logger := utils.GetLogger()

	r := mux.NewRouter()
	r.PathPrefix("/healthz").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		formatted := utils.GetCurrentTime()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		content := `{"status":"OK","time":"` + formatted + `"}`
		w.Write([]byte(content))
		logger.Info("GET", "path", r.RequestURI)
	})

	rootdir := os.Getenv("ROOTDIR")
	if rootdir == "" {
		rootdir = "/tmp/dedup"
	}

	r.PathPrefix("/start").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dryrun := r.URL.Query().Get("dryrun")
		boolVal, err := strconv.ParseBool(dryrun)
		if err != nil {
			boolVal = true
		}
		err = utils.StartProcess(rootdir, boolVal)
		formatted := utils.GetCurrentTime()
		var content string
		if err == nil {
			w.WriteHeader(http.StatusOK)
			content = `{"status":"OK","time":"` + formatted + `"}`
		} else {
			w.WriteHeader(http.StatusBadRequest)
			content = `{"status":"ERROR","time":"` + formatted + `", "error":"` + err.Error() + `"}`
		}
		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(content))
		logger.Info("GET", "start", r.RequestURI)
	})

	r.PathPrefix("/stop").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := utils.StopProcess()
		formatted := utils.GetCurrentTime()
		var content string
		if err == nil {
			w.WriteHeader(http.StatusOK)
			content = `{"status":"OK","time":"` + formatted + `"}`
		} else {
			w.WriteHeader(http.StatusBadRequest)
			content = `{"status":"ERROR","time":"` + formatted + `,"error:"` + err.Error() + `"}`
		}
		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(content))
		logger.Info("GET", "start", r.RequestURI)
	})

	r.PathPrefix("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content := utils.GetStates()
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(content))
		logger.Info("GET", "state", r.RequestURI)
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else {
		port = ":" + port
	}

	err := http.ListenAndServe(port, r)

	if err != nil {
		logger.Error(err.Error())
	}
}
