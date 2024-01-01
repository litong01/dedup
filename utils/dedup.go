package utils

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"time"
)

type dedupState struct {
	InProcess   bool   `json:"inProcess"`
	Duplicates  int    `json:"duplicates"`
	Processed   int    `json:"processed"`
	Unique      int    `json:"unique"`
	CurrentTime string `json:"currentTime"`
	LastError   string `json:"lastError"`
}

var (
	filehash map[hash.Hash]struct{}
	stop     chan bool

	state dedupState

	errStop         = errors.New("stop command received")
	errInProcess    = errors.New("already in process")
	errNotInProcess = errors.New("not in process")
)

func init() {
	filehash = make(map[hash.Hash]struct{})
	stop = make(chan bool)

	state.Duplicates = 0
	state.Processed = 0
	state.InProcess = false
	state.Unique = 0
	state.LastError = ""
}

// GetStates return the dedup state in json format
func GetStates() string {
	state.Unique = len(filehash)
	state.CurrentTime = GetCurrentTime()
	data, _ := json.Marshal(state)
	return string(data)
}

func StartProcess(rootdir string, dryrun bool) error {
	if state.InProcess {
		return errInProcess
	}
	go dedup(rootdir, dryrun)
	return nil
}

func StopProcess() error {
	if !state.InProcess {
		return errNotInProcess
	}
	stop <- true
	return nil
}

func processOneFile(path string, info os.FileInfo, err error, dryrun bool) error {
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
			return nil
		} else {
			logger.Error("Error found", "The error", err)
			return err
		}
	}

	// skip the directory, symbolic link, or socket
	if info.IsDir() || info.Mode()&os.ModeSymlink != 0 || info.Mode()&os.ModeSocket != 0 || info.Mode()&os.ModeType != 0 {
		state.Processed++
		return nil
	}

	thefile, errhere := os.Open(path)
	if errhere != nil {
		logger.Error("Processing file error", "file name", path)
		logger.Error(errhere.Error())
		return err
	}

	// use the sha256 hash as writer to create hash of a file
	thehash := sha256.New()
	if _, errhere = io.Copy(thehash, thefile); errhere != nil {
		logger.Error("Getting hash error", "file name", path)
		logger.Error(errhere.Error())
		return errhere
	}
	// no longer need the file, close it
	thefile.Close()

	// Now if ok true, that means we find a duplicate.
	if _, ok := filehash[thehash]; ok {
		// found duplicate, should remove the file.
		logger.Info("Duplicate instance", "name", path)
		state.Duplicates++
		if !dryrun { // only do the real remove when it is not dryrun
			os.Remove(path)
		}
	} else {
		// not exist in hash, place it in.
		logger.Info("First instance", "name", path)
		filehash[thehash] = struct{}{}
	}
	// finished one file processing, increase the counter
	state.Processed++
	return nil
}

func dedup(rootdir string, dryrun bool) error {
	state.InProcess = true
	state.LastError = ""
	err := filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		select {
		case <-stop: // received the stop signal
			state.InProcess = false
			return errStop
		default: // keep on going
			return processOneFile(path, info, err, dryrun)
		}
	})

	// reach the end, no longer in process
	state.InProcess = false

	// save the error message
	if err != nil {
		state.LastError = err.Error()
	} else {
		state.LastError = ""
	}
	return err
}

func GetCurrentTime() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%07dZ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}
