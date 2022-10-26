package polarisb_syslog_go

import (
	"fmt"
	"github.com/TwiN/go-color"
	"sync"
	"time"
)

type LogEntry struct {
	Time     time.Time
	LogLevel LogLevel
	Event    interface{}
	Message  string
	Actor    string
	MoreInfo map[string]interface{}
	fmt      LogEntryFmt
}

type LogEntryFmt struct {
	NumberOfLinesAfter int
}

type LogSink struct {
	writeToConsole         bool
	writeToChannel         bool
	Channel                *LogSinkChannel
	WriteToConsoleSettings *WriteToConsoleSettings
}

type LogSinkChannel struct {
	SinkChannel chan *LogEntry
	WaitGroup   *sync.WaitGroup
}

type WriteToConsoleSettings struct {
	KeyColor   string
	ValueColor string
}

// SetWriteToConsoleSettings allows setting write to console settings
func (settings *WriteToConsoleSettings) SetWriteToConsoleSettings(newSettings *WriteToConsoleSettings) {
	settings.KeyColor = newSettings.KeyColor
	settings.ValueColor = newSettings.ValueColor
}

// AddPolarisbSystemLog initializes the log sink
func AddPolarisbSystemLog() (logSink *LogSink) {
	newLogSink := &LogSink{
		Channel: &LogSinkChannel{
			SinkChannel: make(chan *LogEntry, 100),
			WaitGroup:   new(sync.WaitGroup),
		},
		writeToChannel: true,
		WriteToConsoleSettings: &WriteToConsoleSettings{
			KeyColor:   color.Reset,
			ValueColor: color.Gray,
		},
	}

	// Add 1 wait for END OF PROGRAM control
	newLogSink.Channel.WaitGroup.Add(1)

	// Configure a Channel auto close goroutine
	go func() {
		newLogSink.Channel.WaitGroup.Wait()
		newLogSink.Channel.CloseChannel()
	}()

	return newLogSink
}

// Close Called at program exit
func (logSink *LogSink) Close() {
	// Indicates that we've reached END OF PROGRAM control
	logSink.Channel.WaitGroup.Done()
	// Now Waiting on old logs
	logSink.Channel.WaitGroup.Wait()
}

// CloseChannel Close sink channel
func (logSink *LogSinkChannel) CloseChannel() {
	close(logSink.SinkChannel)
}

// WriteToConsole Enables writes to the console
func (logSink *LogSink) WriteToConsole() *LogSink {
	logSink.writeToConsole = true
	return logSink
}

// LogThis builds a log entry and sends it to the log sink
func (logSink *LogSink) LogThis(logLevel LogLevel, event interface{}, message string, actor string) {
	logSink.LogThisWithMoreInfo(logLevel, event, message, actor, nil)
}

// LogThisWithMoreInfo builds a log entry with more info
func (logSink *LogSink) LogThisWithMoreInfo(logLevel LogLevel, event interface{}, message string, actor string, moreInfo map[string]interface{}) {
	logSink.LogThisWithMoreInfoAndFmt(logLevel, event, message, actor, moreInfo, LogEntryFmt{NumberOfLinesAfter: 0})
}

// LogThisWithMoreInfoAndFmt builds a log entry with more info and fmt information
func (logSink *LogSink) LogThisWithMoreInfoAndFmt(logLevel LogLevel, event interface{}, message string, actor string, moreInfo map[string]interface{}, entryFmt LogEntryFmt) {
	logEntry := LogEntry{
		Time:     time.Now(),
		LogLevel: logLevel,
		Event:    event,
		Message:  message,
		Actor:    actor,
		MoreInfo: moreInfo,
		fmt:      entryFmt,
	}

	logSink.LogWrite(&logEntry)
}

// LogWrite writes a log entry to the log sink
func (logSink *LogSink) LogWrite(logEntry *LogEntry) {
	if logSink.writeToConsole {
		logSink.consoleOut(logEntry)
	}
	if logSink.writeToChannel {
		logSink.chanOut(logEntry)
	}
}

// consoleOut writes a log entry to the console
func (logSink *LogSink) consoleOut(logEntry *LogEntry) {

	keyColor := logSink.WriteToConsoleSettings.KeyColor
	valueColor := logSink.WriteToConsoleSettings.ValueColor

	fmt.Printf("%s  ", logEntry.Time.Format("2006-01-02 15:04:05"))
	fmt.Printf("%s	", logEntry.LogLevel)
	fmt.Printf(color.Ize(keyColor, "Event:")+color.Ize(valueColor, "\"%s\" "), logEntry.Event)
	fmt.Printf(color.Ize(keyColor, "Message:")+color.Ize(valueColor, "\"%s\" "), logEntry.Message)
	fmt.Printf(color.Ize(keyColor, "Actor:")+color.Ize(valueColor, "\"%s\" "), logEntry.Actor)
	if logEntry.MoreInfo != nil {
		fmt.Printf(color.Ize(keyColor, "More Info:")+color.Ize(valueColor, "\"%s\" "), logEntry.MoreInfo)
	}

	for i := 0; i <= logEntry.fmt.NumberOfLinesAfter; i++ {
		fmt.Printf("\n")
	}

}

// chanOut writes a log entry to the sink channed
func (logSink *LogSink) chanOut(logEntry *LogEntry) {
	logSink.Channel.WaitGroup.Add(1)
	logSink.Channel.SinkChannel <- logEntry
}
