package main

import (
	"github.com/sirupsen/logrus"
	"goGFG/logger"
	"goGFG/storage"
	"sync"
)

var once sync.Once

func main() {
	once.Do(
		func() {
			logger.SetUpLogging()
		})
	storage.GetElastic()
	err := NewHTTPServer()
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"error": err}).Error("Error Starting server")
	}
}
