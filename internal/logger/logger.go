package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Logger struct to hold the log channel and the log file
type Logger struct {
	logChannel chan string
	done       chan bool
}

var (
	instance *Logger
	once     sync.Once
)

const bufferSize = 1000

const LogToFile = false

// NewLogger initializes a new Logger
func newLogger() *Logger {
	l := &Logger{
		logChannel: make(chan string, bufferSize),
		done:       make(chan bool),
	}
	go l.startLogging()
	return l
}

// startLogging listens for log messages and writes them to a file
func (l *Logger) startLogging() {

	if LogToFile {
		logsDir := "logs"

		// Get the current working directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}

		// Create the full path for the logs directory
		logsPath := filepath.Join(currentDir, logsDir)

		// Check if the directory exists
		if _, err := os.Stat(logsPath); os.IsNotExist(err) {
			// Create the logs directory
			err := os.Mkdir(logsPath, 0755) // 0755 is the permission mode
			if err != nil {
				fmt.Println("Error creating logs directory:", err)
				return
			}
		}

		file, err := os.OpenFile("./logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		for {
			select {
			case msg := <-l.logChannel:
				_, err := file.WriteString(fmt.Sprintf("%s: %s\n", time.Now().Format(time.RFC3339), msg))
				if err != nil {
					log.Println("Error writing to log file:", err)
				}
			case <-l.done:
				return
			}
		}
	} else {
		for {
			select {
			case msg := <-l.logChannel:
				// Output log message to stdout
				_, err := fmt.Printf("%s: %s\n", time.Now().Format(time.RFC3339), msg)
				if err != nil {
					log.Println("Error writing to stdout:", err)
				}
			case <-l.done:
				return
			}
		}

	}

}

// Log sends a log message to the log channel
func (l *Logger) Log(message string) {
	l.logChannel <- message
}

// Stop stops the logger
func (l *Logger) Stop() {
	l.done <- true
}

// GetLogger returns the singleton instance of Logger
func GetLogger() *Logger {
	once.Do(func() {
		instance = newLogger()
	})
	return instance
}
