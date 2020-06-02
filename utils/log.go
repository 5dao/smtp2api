package utils

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

//path/logs/xxx_20060102sign.log

var filePrefix string
var logFile *os.File

func init() {
	var err error
	filePrefix, err = getPrefix()
	if err != nil {
		panic(err)
	}

	logCron := cron.New()
	logCron.AddFunc("0 0 0 * * *", MakeDateLog)
	logCron.Start()
}

// MakeDateLog make date.log
func MakeDateLog() {
	defer func() {
		if rev := recover(); rev != nil {
			log.Warn("MakeDateLog recover", rev)
		}
	}()

	now := time.Now()

	var err error

	var fileName string
	fileName = filePrefix + "_" + now.Format("20060102") + ".log"
	oldMask := syscall.Umask(0)
	newLogFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0640)
	syscall.Umask(oldMask)
	if err != nil {
		log.Println("logfile Run os.OpenFile err", err)
		return
	}
	log.SetOutput(newLogFile)

	//close old logfile
	if logFile != nil {
		err = logFile.Close()
		if err != nil {
			log.Println("logfile Run oldLogFile close", err)
		}
	}
	logFile = newLogFile
}

// path/logs/xxx
func getPrefix() (string, error) {
	instanceAbsPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", errors.New("filepath err:" + err.Error())
	}
	instanceName := filepath.Base(instanceAbsPath)

	dir := filepath.Dir(instanceAbsPath)

	oldMask := syscall.Umask(0)
	os.Mkdir(filepath.Join(dir, "logs"), 0750)
	syscall.Umask(oldMask)

	return filepath.Join(dir, "logs", instanceName), nil
}
