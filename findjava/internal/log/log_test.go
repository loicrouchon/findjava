package log

import (
	"findjava/internal/console"
	"reflect"
	"testing"
)

func setTestConsole() *TestConsole {
	_ = SetLogLevel("error")
	testConsole := TestConsole{
		stdout: MessagesHolder{messages: make([]string, 0)},
		stderr: MessagesHolder{messages: make([]string, 0)},
	}
	console.Writer.Stdout = InMemoryWriter{content: &testConsole.stdout}
	console.Writer.Stderr = InMemoryWriter{content: &testConsole.stderr}
	return &testConsole
}

type MessagesHolder struct {
	messages []string
}

type TestConsole struct {
	stdout MessagesHolder
	stderr MessagesHolder
}

func (console *TestConsole) hasMessages(t *testing.T, stdout []string, stderr []string) {
	streamEquals(t, "stdout", console.stdout, stdout)
	streamEquals(t, "stderr", console.stderr, stderr)
}

func streamEquals(t *testing.T, streamName string, stream MessagesHolder, actual []string) {
	if !reflect.DeepEqual(stream.messages, actual) {
		t.Fatalf(`%s == [%#v] but was [%#v]`, streamName, actual, stream.messages)
	}
}

type InMemoryWriter struct {
	content *MessagesHolder
}

func (slice InMemoryWriter) Write(p []byte) (n int, err error) {
	slice.content.messages = append(slice.content.messages, string(p))
	return len(p), nil
}

func (slice InMemoryWriter) get() []string {
	return slice.content.messages
}
