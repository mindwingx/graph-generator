package main

import (
	"github.com/alexflint/go-filemutex"
	command "github.com/mindwingx/graph-generator/cmd/commands"
	"github.com/mindwingx/graph-generator/constants"
	"log"
)

func main() {
	mutex, err := filemutex.New(constants.TmpLockFile)
	if err != nil {
		log.Fatal("/tmp directory does not exist or lock file cannot be created!")
	}

	errLock := mutex.TryLock()
	if errLock != nil {
		log.Fatal("This program is already running on this server!")
	}

	err = mutex.Lock()
	if err != nil {
		log.Fatal(err)
	}

	// initialize root commands and run related commands
	command.Exec()

	err = mutex.Unlock()
	if err != nil {
		log.Fatal(err)
	}
}
