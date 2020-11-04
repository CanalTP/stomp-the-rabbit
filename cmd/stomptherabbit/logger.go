package main

import (
	"os"
	"runtime"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// GetLogger initializes and return a logger
func getLogger(version string, json bool) *log.Entry {
	println(json)
	if json {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	log.SetOutput(os.Stdout)

	log.SetLevel(log.DebugLevel)
	logger := log.WithFields(log.Fields{
		"version": version,
		"runtime": runtime.Version(),
	})

	return logger
}
