package logger

import (
	"github.com/jhidalgoesp/example-services/tests"
	"testing"
)

func TestInitLogger(t *testing.T) {
	t.Run("Logger is not nil", func(t *testing.T) {
		logger, _ := InitLogger("test-tag")
		tests.AssertKindIsPointer(t, logger)
	})
}
