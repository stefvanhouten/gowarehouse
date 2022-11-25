package logger

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	DefaultLogger *logrus.Logger
)

type WriterHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *WriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
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
	DefaultLogger.SetFormatter(&logrus.JSONFormatter{})

	f, err := os.OpenFile(os.Getenv("LOGDIR")+"log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err.Error())
	}

	if env := os.Getenv("ENVIRONMENT"); env == "dev" {
		// In development environments we want DEBUG and INFO messages in the console,
		// INFO and above can go into the regular logfile.
		DefaultLogger.SetOutput(ioutil.Discard)
		DefaultLogger.SetLevel(logrus.DebugLevel)
		DefaultLogger.AddHook(&WriterHook{
			Writer: os.Stderr,
			LogLevels: []logrus.Level{
				logrus.DebugLevel,
				logrus.InfoLevel,
			},
		})

		DefaultLogger.AddHook(&WriterHook{
			Writer: f,
			LogLevels: []logrus.Level{
				logrus.WarnLevel,
				logrus.PanicLevel,
				logrus.ErrorLevel,
			},
		})
	} else {
		DefaultLogger.SetLevel(logrus.InfoLevel)
		DefaultLogger.SetOutput(f)
	}
}
