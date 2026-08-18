package logger

import (
	"bytes"
	"io"
	"sync"
	"testing"
)

type mockLogger struct {
	mu   sync.Mutex
	logs []*Message
}

func (m *mockLogger) Log(msg *Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	lineCopy := make([]byte, len(msg.Line))
	copy(lineCopy, msg.Line)
	m.logs = append(m.logs, &Message{
		Source:    msg.Source,
		Line:      lineCopy,
		Timestamp: msg.Timestamp,
	})
	return nil
}

func TestCopierPartialLine(t *testing.T) {
	src := bytes.NewBufferString("final-partial-line")
	logger := &mockLogger{}
	copier := NewCopier(map[string]io.Reader{"stdout": src}, logger)
	copier.Run()

	logger.mu.Lock()
	defer logger.mu.Unlock()

	if len(logger.logs) != 1 {
		t.Fatalf("expected 1 log message, got %d", len(logger.logs))
	}

	expected := "final-partial-line"
	actual := string(logger.logs[0].Line)
	if actual != expected {
		t.Errorf("expected log line %q, got %q", expected, actual)
	}
}
