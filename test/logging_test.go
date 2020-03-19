package test

import (
	"testing"

	"github.com/dashjay/logging"
)

func TestLogging(t *testing.T) {
	logging.Info("info test")
	logging.Debug("debug test")
	logging.Error("error test")
	// logging.Fatal("fatal test")
	logging.Warn("warn test")
}
