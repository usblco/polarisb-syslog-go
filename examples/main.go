package main

import (
	"fmt"
	"github.com/TwiN/go-color"
	polarisb "github.com/usblco/polarisb-syslog-go"
)

func clientChannelListener(polarisbLogging *polarisb.LogSink) {
	for logEntry := range polarisbLogging.Channel.SinkChannel {
		if logEntry != nil {
			fmt.Printf("Event Received | Event: %s | Message: %s | Actor: %s | MoreInfo: %s \n", logEntry.Event, logEntry.Message, logEntry.Actor, logEntry.MoreInfo)
			// Called to indicate that log entry has been processed
			polarisbLogging.Channel.WaitGroup.Done()
		}
	}
}

func main() {
	polarisbLogging := polarisb.AddPolarisbSystemLog()
	polarisbLogging.WriteToConsoleSettings.SetWriteToConsoleSettings(&polarisb.WriteToConsoleSettings{
		KeyColor:   color.Blue,
		ValueColor: color.Green,
	})

	go clientChannelListener(polarisbLogging)

	go bigLogger(polarisbLogging)

	polarisbLogging.LogThis(
		polarisb.Information,
		"Hello World 1",
		"Polarisb Service Manager Database Initialized",
		"system")

	polarisbLogging.LogThisWithMoreInfo(
		polarisb.Information,
		"Hello World 2",
		"Service Manager Network Initialized",
		"system",
		map[string]interface{}{
			"test":  "test",
			"test2": "test2",
		})

	polarisbLogging.LogThisWithMoreInfoAndFmt(
		polarisb.Information,
		"Hello World 3",
		"Service Manager Network Initialized",
		"system",
		map[string]interface{}{
			"test":  "test",
			"test2": "test2",
		},
		polarisb.LogEntryFmt{NumberOfLinesAfter: 0})

	polarisbLogging.LogThis(
		polarisb.Information,
		"Hello World 4",
		"Polarisb Service Manager Database Initialized",
		"system")

	fmt.Println("Waiting for log entries to be processed hit end of program")

	polarisbLogging.Close()

}

func bigLogger(polarisbLogging *polarisb.LogSink) {
	for i := 0; i < 1000; i++ {
		polarisbLogging.LogThis(
			polarisb.Information,
			fmt.Sprintf("Hello World %d", i),
			"Polarisb Service Manager Database Initialized",
			"system")
	}
}
