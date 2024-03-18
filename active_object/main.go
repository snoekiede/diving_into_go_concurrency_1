package main

import (
	"fmt"
	"sync"
	"time"
)

type LogSeverityType string

const (
	Info    LogSeverityType = "Info"
	Warning LogSeverityType = "Warning"
	Error   LogSeverityType = "Error"
)

type LogMessage struct {
	Severity LogSeverityType
	Message  string
}

type ActiveLogger struct {
	logQueue  chan *LogMessage
	stopEvent chan struct{}
	wg        sync.WaitGroup
}

func createActiveLogger() *ActiveLogger {
	return &ActiveLogger{
		logQueue:  make(chan *LogMessage, 10),
		stopEvent: make(chan struct{}),
	}
}

func (l *ActiveLogger) Log(message *LogMessage) {
	l.logQueue <- message
}

func (l *ActiveLogger) StartLogProcessor() {
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		for {
			select {
			case message := <-l.logQueue:
				l.processMessage(message)
			case <-l.stopEvent:
				return
			}
		}
	}()
}

func (l *ActiveLogger) StopLogProcessor() {
	close(l.stopEvent)
	l.wg.Wait()
}

func (l *ActiveLogger) processMessage(m *LogMessage) {
	fmt.Printf("Processing: (%s): %s\n", m.Severity, m.Message)
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("Processed: (%s): %s\n", m.Severity, m.Message)
}

func main() {
	logger := createActiveLogger()

	logger.StartLogProcessor()
	for i := 1; i <= 30; i++ {
		message := &LogMessage{
			Message:  fmt.Sprintf("Message number is %d", i),
			Severity: Info,
		}
		logger.Log(message)
	}
	message := &LogMessage{
		Message:  "There has been an error",
		Severity: Error,
	}
	logger.Log(message)

	time.Sleep(2 * time.Second)
	logger.StopLogProcessor()
	fmt.Println("Stopped log processor")
}
