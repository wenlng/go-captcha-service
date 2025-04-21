package helper

type OutputLogType string

const (
	OutputLogTypeWarn  OutputLogType = "warn"
	OutputLogTypeInfo                = "info"
	OutputLogTypeError               = "error"
	OutputLogTypeDebug               = "debug"
)

type OutputLogCallback = func(logType OutputLogType, message string)
