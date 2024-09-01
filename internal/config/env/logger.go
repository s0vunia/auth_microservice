package env

import (
	"errors"
	"os"
	"strconv"

	"github.com/s0vunia/auth_microservice/internal/config"
)

const (
	filenameEnvName       = "LOG_FILENAME"
	fileMaxSizeEnvName    = "LOG_FILE_MAX_SIZE"
	fileMaxBackupsEnvName = "LOG_FILE_MAX_BACKUPS"
	fileMaxAgeEnvName     = "LOG_FILE_MAX_AGE"
	logLevelEnvName       = "LOG_LEVEL"
)

type loggerConfig struct {
	filename       string
	fileMaxSize    int
	fileMaxBackups int
	fileMaxAge     int
	logLevel       string
}

// NewLoggerConfig - creates new logger config
func NewLoggerConfig() (config.LoggerConfig, error) {
	parseInt := func(s string) int {
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return i
	}

	filename := os.Getenv(filenameEnvName)
	if len(filename) == 0 {
		return nil, errors.New("filename not found")
	}

	fileMaxSize := os.Getenv(fileMaxSizeEnvName)
	if len(fileMaxSize) == 0 {
		return nil, errors.New("fileMaxSize not found")
	}

	fileMaxBackups := os.Getenv(fileMaxBackupsEnvName)
	if len(fileMaxBackups) == 0 {
		return nil, errors.New("fileMaxBackups not found")
	}

	fileMaxAge := os.Getenv(fileMaxAgeEnvName)
	if len(fileMaxAge) == 0 {
		return nil, errors.New("fileMaxAge not found")
	}

	logLevel := os.Getenv(logLevelEnvName)
	if len(logLevel) == 0 {
		return nil, errors.New("logLevel not found")
	}

	return &loggerConfig{
		filename:       filename,
		fileMaxSize:    parseInt(fileMaxSize),
		fileMaxBackups: parseInt(fileMaxBackups),
		fileMaxAge:     parseInt(fileMaxAge),
		logLevel:       logLevel,
	}, nil
}

func (l *loggerConfig) FileName() string {
	return l.filename
}

func (l *loggerConfig) Level() string {
	return l.logLevel
}

func (l *loggerConfig) MaxSize() int {
	return l.fileMaxSize
}

func (l *loggerConfig) MaxBackups() int {
	return l.fileMaxBackups
}

func (l *loggerConfig) MaxAge() int {
	return l.fileMaxAge
}
