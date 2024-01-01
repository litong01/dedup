package utils

import (
	"crypto/sha256"
	"hash"
	"io"
	"os"
	"path/filepath"

	goslog "golang.org/x/exp/slog"
)

var (
	filehash      map[hash.Hash]struct{}
	fileprocessed int
	Logger        *goslog.Logger
)

func init() {
	filehash = make(map[hash.Hash]struct{})
	fileprocessed = 0
}

func GetStates() (int, int) {
	return fileprocessed, len(filehash)
}

func Dedup(rootdir string) error {
	err := filepath.Walk(rootdir, func(path string, info os.FileInfo, err error) error {
		fileprocessed++
		if err != nil {
			return err
		}
		if info.IsDir() {
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
			Logger.Info("Duplicate instance", "name", path)
			os.Remove(path)
		} else {
			// not exist in hash, place it in.
			Logger.Info("First instance", "name", path)
			filehash[thehash] = struct{}{}
		}

		return nil
	})

	return err
}
