package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type LoggerConfig interface {
	GetEnvironment() string
	GetLogDir() string
}

type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
}

const (
	DEVELOPMENT = "DEV"
	PRODUCTION  = "PROD"
)

// Global logger instance that can be used throughout the application with the
// predefined configuration.
var (
	DefaultLogger *logrus.Logger
)

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *WriterHook) Levels() []logrus.Level {
	return hook.LogLevels
}

// Setup the logger for development mode, logging to stdout and file based on LogLevel.
func setupDebugLogger(logfile *os.File) {
	// This removes all configured output hooks, allowing us to overwrite
	// the default output location.
	DefaultLogger.SetOutput(ioutil.Discard)

	// In development environments we want DEBUG and INFO messages in the console,
	// INFO and above can go into the regular logfile.
	DefaultLogger.SetLevel(logrus.DebugLevel)

	// DEBUG and INFO to console
	DefaultLogger.AddHook(&WriterHook{
		Writer: os.Stderr,
		LogLevels: []logrus.Level{
			logrus.DebugLevel,
			logrus.InfoLevel,
		},
		// The logs that go into the console should be human readable and have a nice
		// color and format.
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   false,
			ForceColors:     true,
		},
	})

	// Everything else can go into the logfile.
	DefaultLogger.AddHook(&WriterHook{
		Writer: logfile,
		LogLevels: []logrus.Level{
			logrus.WarnLevel,
			logrus.PanicLevel,
			logrus.ErrorLevel,
			logrus.FatalLevel,
			logrus.TraceLevel,
		},
		Formatter: &logrus.JSONFormatter{},
	})
}

// Setup the logger for production mode. This means logging to a file only.
func setupProductionLogger(logfile *os.File) {
	DefaultLogger.SetFormatter(&logrus.JSONFormatter{})

	// In production all INFO and above logs will be written to file.
	DefaultLogger.SetLevel(logrus.InfoLevel)
	DefaultLogger.SetOutput(logfile)
}

// Setup the logger based on passed config.
// This funcion will create a new DefaultLogger instance, and only has to be called once.
func Setup(conf LoggerConfig) {
	// TODO: Find out the go way of setting up a logger instance like this.
	// Singleton pattern?
	DefaultLogger = logrus.New()

	// Open the logfile.
	logfile, err := os.OpenFile(
		conf.GetLogDir()+"log.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644,
	)

	if err != nil {
		panic(err.Error())
	}

	if strings.ToUpper(conf.GetEnvironment()) == PRODUCTION {
		setupProductionLogger(logfile)
	} else {
		setupDebugLogger(logfile)
	}
}
