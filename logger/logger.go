package logger

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// Global logger instance that can be used throughout the application with the
// predefined configuration.
var (
	DefaultLogger *logrus.Logger
)

type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
}

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

func init() {
	DefaultLogger = logrus.New()
	// Open the logfile.
	f, err := os.OpenFile(
		os.Getenv("LOGDIR")+"log.log",
		os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644,
	)

	if err != nil {
		panic(err.Error())
	}

	if env := os.Getenv("ENVIRONMENT"); env != "dev" {
		DefaultLogger.SetFormatter(&logrus.JSONFormatter{})

		// In production all INFO and above logs will be written to file.
		DefaultLogger.SetLevel(logrus.InfoLevel)
		DefaultLogger.SetOutput(f)
		return
	}

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
		Formatter: &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   false,
			ForceColors:     true,
		},
	})

	// Everything else can go into the logfile.
	DefaultLogger.AddHook(&WriterHook{
		Writer: f,
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
