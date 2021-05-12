// Package logger is the logging infrastructure of the project
package logger

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

var log logr.Logger
var zlog *zap.Logger

// Setup initialize the logging infrastructure using production mode if needed
func Setup(production bool) {
	var err error

	if production {
		zlog, err = zap.NewProduction()
	} else {
		zlog, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}

	log = zapr.NewLogger(zlog)
}

// Log returns the root logger
func Log() logr.Logger {
	return log
}
