package utils

import (
    "log"
    "os"

    "github.com/sirupsen/logrus"
)

// SetupLogger настраивает логгер
func SetupLogger(debug bool) {
    logger := logrus.New()
    logger.SetOutput(os.Stdout)

    if debug {
        logger.SetLevel(logrus.DebugLevel)
        logger.Info("Режим отладки включен")
    } else {
        logger.SetLevel(logrus.InfoLevel)
    }

    log.Logger = logger
}