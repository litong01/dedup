package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	goslog "golang.org/x/exp/slog"

	"github.com/dedup/utils"
	"github.com/gorilla/mux"
)

var (
	Logger *goslog.Logger
	opts   goslog.HandlerOptions
)

func init() {
	doLog := os.Getenv("DOLOG")
	// TODO getting configuration parameters of the control,
	// then use these parameters to customize the logger.
	if doLog == "" {
		opts.Level = goslog.LevelError
	} else {
		opts.Level = goslog.LevelInfo
	}
	Logger = goslog.New(goslog.NewJSONHandler(os.Stdout, &opts))
	goslog.SetDefault(Logger)
}

func main() {

	r := mux.NewRouter()
	r.PathPrefix("/healthz").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		formatted := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%07dZ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		content := `{"status":"OK","time":"` + formatted + "\"}"
		w.Write([]byte(content))
		Logger.Info("GET", "path", r.RequestURI)
	})

	r.PathPrefix("/").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.RequestURI))
		Logger.Info("GET", "path", r.RequestURI)
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = ":8080"
	} else {
		port = ":" + port
	}

	rootdir := os.Getenv("ROOTDIR")

	go func() {
		Logger.Info(rootdir)
		utils.Dedup(rootdir)
	}()

	err := http.ListenAndServe(port, r)

	if err != nil {
		Logger.Error(err.Error())
	}
}
