package logger_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/logger"
)

// parseLog is a small helper that takes the raw JSON log output from zerolog
// and unmarshals it into a Go map so we can assert on structured fields.
//
// This is better than checking for raw substrings, because zerolog always
// outputs JSON logs, and relying on exact formatting or timestamps would
// make the tests brittle.
func parseLog(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	var logEntry map[string]any
	err := json.Unmarshal(buf.Bytes(), &logEntry)
	require.NoError(t, err)

	return logEntry
}

// TestNewZerologLogger_DefaultLevel ensures that when an invalid log level
// is provided, the logger defaults to "info".
func TestNewZerologLogger_DefaultLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("invalid-level", &buf)

	l.Info("hello %s", "world")

	entry := parseLog(t, &buf)
	assert.Equal(t, "info", entry["level"])
	assert.Equal(t, "hello world", entry["message"])
}

// TestNewZerologLogger_DebugLevel ensures that when the level is set to "debug",
// debug messages are logged correctly.
func TestNewZerologLogger_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("debug", &buf)

	l.Debug("debug %d", 123)

	entry := parseLog(t, &buf)
	assert.Equal(t, "debug", entry["level"])
	assert.Equal(t, "debug 123", entry["message"])
}

// TestLogger_Info verifies that Info messages appear with the correct level and message.
func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("info", &buf)

	l.Info("info message %d", 1)

	entry := parseLog(t, &buf)
	assert.Equal(t, "info", entry["level"])
	assert.Equal(t, "info message 1", entry["message"])
}

// TestLogger_Warn verifies that Warn messages appear with the correct level and message.
func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("warn", &buf)

	l.Warn("warn message %s", "X")

	entry := parseLog(t, &buf)
	assert.Equal(t, "warn", entry["level"])
	assert.Equal(t, "warn message X", entry["message"])
}

// TestLogger_Debug verifies that Debug messages appear with the correct level and message.
func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("debug", &buf)

	l.Debug("debug message %d", 42)

	entry := parseLog(t, &buf)
	assert.Equal(t, "debug", entry["level"])
	assert.Equal(t, "debug message 42", entry["message"])
}

// TestLogger_Error verifies that Error messages appear with the correct level and message.
func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("error", &buf)

	l.Error("error message %s", "oops")

	entry := parseLog(t, &buf)
	assert.Equal(t, "error", entry["level"])
	assert.Equal(t, "error message oops", entry["message"])
}

// TestLogger_Panic verifies that Panic messages are logged AND that the method panics.
func TestLogger_Panic(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("panic", &buf)

	assert.Panics(t, func() {
		l.Panic("panic message %s", "boom")
	})

	entry := parseLog(t, &buf)
	assert.Equal(t, "panic", entry["level"])
	assert.Equal(t, "panic message boom", entry["message"])
}

// We deliberately skip Fatal because it calls os.Exit(1)
// and would stop the entire test process.
