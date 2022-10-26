package polarisb_syslog_go

type LogLevel int64

const (
	Debug       LogLevel = iota
	Information          // Use this for normal logging
	Warning              // Use this for logging that is not an error, but may be a problem
	Error                // Use this for logging that is an error
	Fatal                // Use this for logging that is a fatal error

	Success  // Use this for logging that is a success
	Failure  // Use this for logging that is a failure
	Critical // Use this for logging that is a critical failure, but can also be a warning
)

func (level LogLevel) String() string {
	switch level {
	case Debug:
		return "Debug"
	case Information:
		return "Information"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Fatal:
		return "Fatal"
	case Success:
		return "Success"
	case Failure:
		return "Failure"
	case Critical:
		return "Critical"
	}

	return "Unknown"
}
