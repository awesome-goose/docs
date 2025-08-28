<!--

Handlers
Log to files and syslog
StreamHandler: Logs records into any PHP stream, use this for log files.
RotatingFileHandler: Logs records to a file and creates one log file per day. It will also delete files older than $maxFiles. You should use logrotate for high profile setups though, this is just meant as a quick and dirty solution.
SyslogHandler: Logs records to the syslog.

Send alerts and emails
NativeMailerHandler: Sends emails using PHP's mail() function.
SlackHandler: Logs records to a Slack account using the Slack API (complex setup).

Log to databases
MongoDBHandler: Handler to write records in MongoDB via a Mongo extension connection.

Wrappers / Special Handlers
GroupHandler: This handler groups other handlers. Every record received is sent to all the handlers it is configured with.
FilterHandler: This handler only lets records of the given levels through to the wrapped handler.
SamplingHandler: Wraps around another handler and lets you sample records if you only want to store some of them.
NoopHandler: This handler handles anything by doing nothing. It does not stop processing the rest of the stack. This can be used for testing, or to disable a handler when overriding a configuration.

Formatters
LineFormatter: Formats a log record into a one-line string.
JsonFormatter: Encodes a log record into json.
SyslogFormatter: Used to format log records in RFC 5424 / syslog format. This can be used to output a syslog-style file that can then be consumed by tools like lnav.

Processors
IntrospectionProcessor: Adds the line/file/class/method from which the log call originated.
WebProcessor: Adds the current request URI, request method and client IP to a log record.
MemoryUsageProcessor: Adds the current memory usage to a log record.
MemoryPeakUsageProcessor: Adds the peak memory usage to a log record.
ProcessIdProcessor: Adds the process id to a log record.
UidProcessor: Adds a unique identifier to a log record.
TagProcessor: Adds an array of predefined tags to a log record.
HostnameProcessor: Adds the current hostname to a log record. -->

<!--
package log

// Awesome Goose
// Log
//
// Log is a simple logging package that provides a unified interface for logging messages.
// It supports different log levels and allows you to set a custom logger.
//
// TODO:
// as external package
// processors/stream.go
// processors/slack.go
// processors/mongo.go
// processors/mail.go
// processors/http.go
//
// USAGE:
// - add channel + loggers to platform config
// - platform does *log.Log.Add(channel, loggers)
// - services/other consumers can use:
// - - *log.Log.Use
// - - *log.Contracts.ErrorLogLevel
// - - *log.Log.\*\*\*\*
//
// Advance:
// - create modifiers
// - create formatters
// - create processor
// - connect/mix them as follows
// ```
// myModifier = NewMyModifier()
// myFormatter = NewMyFormatter()
// myProcessor = NewMyProcessor()

// myChannel = "myChannel"
// myLoggers = []\*log.logger{
// log.NewLogger(myModifier, myFormatter, myProcessor),
// ... // add more loggers if needed
// }

// Log = log.Log.Add(myChannel, myLoggers...)
//
// in services/consumer, use: log.Log.Use(myChannel).Debug(...)
// ```

LogLevel -->
