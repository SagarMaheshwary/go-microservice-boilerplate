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

	l.Info("hello world", logger.Field{Key: "foo", Value: "bar"})

	entry := parseLog(t, &buf)
	assert.Equal(t, "info", entry["level"])
	assert.Equal(t, "hello world", entry["message"])
	assert.Equal(t, "bar", entry["foo"])
}

// TestNewZerologLogger_DebugLevel ensures that when the level is set to "debug",
// debug messages are logged correctly.
func TestNewZerologLogger_DebugLevel(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("debug", &buf)

	l.Debug("debugging", logger.Field{Key: "count", Value: 123})

	entry := parseLog(t, &buf)
	assert.Equal(t, "debug", entry["level"])
	assert.Equal(t, "debugging", entry["message"])
	assert.Equal(t, float64(123), entry["count"]) // JSON numbers decode as float64
}

// TestLogger_Info verifies that Info messages appear with the correct level and message.
func TestLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("info", &buf)

	l.Info("info message", logger.Field{Key: "id", Value: 1})

	entry := parseLog(t, &buf)
	assert.Equal(t, "info", entry["level"])
	assert.Equal(t, "info message", entry["message"])
	assert.Equal(t, float64(1), entry["id"])
}

// TestLogger_Warn verifies that Warn messages appear with the correct level and message.
func TestLogger_Warn(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("warn", &buf)

	l.Warn("warn message", logger.Field{Key: "tag", Value: "X"})

	entry := parseLog(t, &buf)
	assert.Equal(t, "warn", entry["level"])
	assert.Equal(t, "warn message", entry["message"])
	assert.Equal(t, "X", entry["tag"])
}

// TestLogger_Debug verifies that Debug messages appear with the correct level and message.
func TestLogger_Debug(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("debug", &buf)

	l.Debug("debugging", logger.Field{Key: "num", Value: 42})

	entry := parseLog(t, &buf)
	assert.Equal(t, "debug", entry["level"])
	assert.Equal(t, "debugging", entry["message"])
	assert.Equal(t, float64(42), entry["num"])
}

// TestLogger_Error verifies that Error messages appear with the correct level and message.
func TestLogger_Error(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("error", &buf)

	l.Error("error happened", logger.Field{Key: "err", Value: "oops"})

	entry := parseLog(t, &buf)
	assert.Equal(t, "error", entry["level"])
	assert.Equal(t, "error happened", entry["message"])
	assert.Equal(t, "oops", entry["err"])
}

// TestLogger_Panic verifies that Panic messages are logged AND that the method panics.
func TestLogger_Panic(t *testing.T) {
	var buf bytes.Buffer
	l := logger.NewZerologLogger("panic", &buf)

	assert.Panics(t, func() {
		l.Panic("panic message", logger.Field{Key: "boom", Value: true})
	})

	entry := parseLog(t, &buf)
	assert.Equal(t, "panic", entry["level"])
	assert.Equal(t, "panic message", entry["message"])
	assert.Equal(t, true, entry["boom"])
}
