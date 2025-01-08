package logger

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	logBuffer []map[string]interface{}
	mu        sync.Mutex
)

func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

func LogEntry(entry map[string]interface{}) {
	mu.Lock()
	defer mu.Unlock()
	logBuffer = append(logBuffer, entry)
}

func FlushLogs(filename string) error {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(logBuffer)
}
