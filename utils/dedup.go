package utils

import (
	"crypto/sha256"
	"errors"
	"hash"
	"io"
	"os"
	"path/filepath"
)

var (
	filehash      map[hash.Hash]struct{}
	fileprocessed int
	dupfound      int
	stop          chan bool
	inProcess     bool

	errStop         = errors.New("stop command received")
	errInProcess    = errors.New("already in process")
	errNotInProcess = errors.New("not in process")
)

func init() {
	filehash = make(map[hash.Hash]struct{})
	fileprocessed = 0
	dupfound = 0
	stop = make(chan bool)
	inProcess = false
}

func GetStates() (int, int) {
	return fileprocessed, len(filehash)
}

func StartProcess(rootdir string, dryrun bool) error {
	if inProcess {
		return errInProcess
	}
	go dedup(rootdir, dryrun)
	return nil
}

func StopProcess() error {
	if !inProcess {
		return errNotInProcess
	}
	stop <- true
	return nil
}

func processOneFile(path string, info os.FileInfo, err error, dryrun bool) error {
	if err != nil {
		return err
	}
	// skip the directory
	if info.IsDir() {
		fileprocessed++
		return nil
	}
	thefile, err := os.Open(path)
	if err != nil {
		return err
	}

	// use the sha256 hash as writer to create hash of a file
	thehash := sha256.New()
	if _, err = io.Copy(thehash, thefile); err != nil {
		return err
	}
	// no longer need the file, close it
	thefile.Close()

	// Now if ok true, that means we find a duplicate.
	if _, ok := filehash[thehash]; ok {
		// found duplicate, should remove the file.
		logger.Info("Duplicate instance", "name", path)
		dupfound++
		if !dryrun { // only do the real remove when it is not dryrun
			os.Remove(path)
		}
	} else {
		// not exist in hash, place it in.
		logger.Info("First instance", "name", path)
		filehash[thehash] = struct{}{}
	}
	// finished one file processing, increase the counter
	fileprocessed++
	return nil
}

func dedup(rootdir string, dryrun bool) error {
	inProcess = true
	err := filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		select {
		case <-stop: // received the stop signal
			inProcess = false
			return errStop
		default: // keep on going
			return processOneFile(path, info, err, dryrun)
		}
	})

	// reach the end, no longer in process
	inProcess = false
	return err
}
