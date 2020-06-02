package utils

import (
	"log"
	"os"
	"path/filepath"

	"time"
)

// make log file date

var logNow, logTomorrow time.Time
var LogFile *os.File

// /xx/xx/logs/xxx +  _20060102.log
var logAbsPathPrefix string

func init() {
	if !logGetFilePrefix() {
		return
	}

	if logAbsPathPrefix == "" {
		return
	}

	go logDateFile()
}

func logDateFile() {
	defer func() {
		if rev := recover(); rev != nil {
			go logDateFile()
		}
	}()

	logMakeFile()

	for {
		select {
		case <-time.After(logTomorrow.Sub(logNow)):
			logMakeFile()
		}
	}
}

func logMakeFile() {
	logNow = time.Now()
	logTomorrow = time.Date(logNow.Year(), logNow.Month(), logNow.Day()+1, 0, 0, 0, 0, logNow.Location())
	if LogFile != nil {
		LogFile.Close()
	}
	fileName := logAbsPathPrefix + "_" + logNow.Format("20060102") + ".log"
	LogFile, _ = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, os.ModePerm)
	log.SetOutput(LogFile)
}

// xx/bin/exe  xx/logs
func logGetFilePrefix() bool {
	instanceAbsPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Println("logGetFileName err", err)
		return false
	}
	instanceName := filepath.Base(instanceAbsPath)

	//
	excDir := filepath.Dir(instanceAbsPath)

	logsDir := filepath.Join(excDir, "logs")

	_, err = os.Stat(logsDir)
	if err != nil {

		err2 := os.Mkdir(logsDir, os.ModePerm)
		if err2 != nil {
			log.Println("logGetFilePrefix mkdir err", err2)
			return false
		}
	}

	logAbsPathPrefix = filepath.Join(logsDir, instanceName)

	return true
}
