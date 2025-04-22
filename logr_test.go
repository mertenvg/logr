package logr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var buf = &bytes.Buffer{}
var jsb = &bytes.Buffer{}
var errb = &bytes.Buffer{}

func TestMain(m *testing.M) {
	SetMeta(MetaData{"application": "logr"})
	AddWriter(buf)
	AddWriter(jsb, WithFormatter(FormatJSON))
	AddWriter(errb, WithFilter(Critical))
	// AddWriter(os.Stdout, WithFormatter(FormatWithColours))
	// AddWriter(os.Stdout, WithFormatter(FormatJSON))

	code := m.Run()

	buf.Reset()
	jsb.Reset()

	os.Exit(code)
}

func mustReadBuffer(buf *bytes.Buffer, t *testing.T) []byte {
	bs, err := ioutil.ReadAll(buf)
	if err != nil {
		t.Errorf("failed to read buffer: %v", err)
	}
	return bs
}

func TestAddWriter(t *testing.T) {
	awbuf := &bytes.Buffer{}
	stop := AddWriter(awbuf)

	Info("TestAddWriter message 1")
	Wait()

	b := mustReadBuffer(awbuf, t)
	if !bytes.Contains(b, []byte("TestAddWriter message 1")) {
		t.Error("expected buffer to contain 'TestAddWriter message 1'")
	}

	stop()

	Info("TestAddWriter message 2")
	Wait()

	b2 := mustReadBuffer(awbuf, t)
	if len(b2) > 0 {
		t.Error("expected buffer to be empty")
	}
}

func testMessageFunc(t *testing.T, ft Type, f func(v ...any) string, args ...any) {
	f(args...)
	Wait()

	b := mustReadBuffer(buf, t)
	if !bytes.Contains(b, []byte("| "+ft.Rune()+" |")) {
		t.Errorf("expected buffer to contain '| \" + ft.Rune() + \" |'. Got: %s", b)
	}
	if !bytes.Contains(b, []byte(strings.TrimSpace(fmt.Sprintln(args...)))) {
		t.Errorf("expected buffer to contain '%v'. Got: %s", strings.TrimSpace(fmt.Sprintln(args...)), b)
	}

	j := mustReadBuffer(jsb, t)
	if !bytes.Contains(j, []byte("\"type\":\""+ft.String()+"\"")) {
		t.Errorf("expected json buffer result to contain '\"type\":\". Got: %s"+ft.String()+"\"'", j)
	}
	if !bytes.Contains(j, []byte(strings.TrimSpace(fmt.Sprintln(args...)))) {
		t.Errorf("expected json buffer to contain '%v'. Got: %s", strings.TrimSpace(fmt.Sprintln(args...)), j)
	}

	e := mustReadBuffer(errb, t)
	if len(e) > 0 {
		if ft&Critical != ft {
			t.Errorf("expected error buffer to be empty. Got: %s", e)
		}
	}
}

func testMessagefFunc(t *testing.T, ft Type, f func(msg string, v ...any) string, msg string, args ...any) {
	f(msg, args...)
	Wait()

	b := mustReadBuffer(buf, t)
	if !bytes.Contains(b, []byte("| "+ft.Rune()+" |")) {
		t.Errorf("expected buffer to contain '| \" + ft.Rune() + \" |'. Got: %s", b)
	}
	if !bytes.Contains(b, []byte(fmt.Sprintf(msg, args...))) {
		t.Errorf("expected buffer to contain '%v'. Got: %s", fmt.Sprintf(msg, args...), b)
	}

	j := mustReadBuffer(jsb, t)
	if !bytes.Contains(j, []byte("\"type\":\""+ft.String()+"\"")) {
		t.Errorf("expected json buffer result to contain '\"type\":\". Got: %s"+ft.String()+"\"'", j)
	}
	if !bytes.Contains(j, []byte(fmt.Sprintf(msg, args...))) {
		t.Errorf("expected json buffer to contain '%v'. Got: %s", fmt.Sprintf(msg, args...), j)
	}

	e := mustReadBuffer(errb, t)
	if len(e) > 0 {
		if ft&Critical != ft {
			t.Errorf("expected error buffer to be empty. Got: %s", e)
		}
	}
}

func TestSuccess(t *testing.T) {
	testMessageFunc(t, S, Success, "Test success message", "with argument concatenation")
}

func TestSuccessf(t *testing.T) {
	testMessagefFunc(t, S, Successf, "Test success message with replacement '%v'", "value")
}

func TestInfo(t *testing.T) {
	testMessageFunc(t, I, Info, "Test info message", "with argument concatenation")
}

func TestInfof(t *testing.T) {
	testMessagefFunc(t, I, Infof, "Test info message with replacement '%v'", "value")
}

func TestDebug(t *testing.T) {
	testMessageFunc(t, D, Debug, "Test debug message", "with argument concatenation")
}

func TestDebugf(t *testing.T) {
	testMessagefFunc(t, D, Debugf, "Test debug message with replacement '%v'", "value")
}

func TestError(t *testing.T) {
	testMessageFunc(t, E, Error, "Test error message", "with argument concatenation")
}

func TestErrorf(t *testing.T) {
	testMessagefFunc(t, E, Errorf, "Test error message with replacement '%v'", "value")
}

func TestWith(t *testing.T) {
	With(Meta{"test": "TestWith"}).Error("Test error message with additional meta data")
	Wait()

	b := mustReadBuffer(buf, t)
	if !bytes.Contains(b, []byte("test:TestWith")) {
		t.Errorf("expected buffer to contain 'test:TestWith'. Got: %s", b)
	}
}

func ExampleAddWriter() {
	AddWriter(os.Stdout)
	// which is equivalent to
	// 		AddWriter(os.Stdout, WithFormatter(FormatDefault), WithFilter(All))
}

func ExampleFormatWithColours() {
	AddWriter(os.Stdout, WithFormatter(FormatWithColours))
}

func ExampleWithFilter() {
	// None Type = iota      // None
	// P    Type = 1 << iota // Panic
	// E                     // Error
	// W                     // Warning
	// I                     // Info
	// D                     // Debug
	// S                     // Success
	// Critical = P | E           // Panic and Error
	// Monitor  = Critical | W    // Panic, Error, and Warning
	// Verbose  = Monitor | I | S // Panic, Error, Warning, Info, and Success
	// All      = Verbose | D     // Panic, Error, Warning, Info, Success, and Debug
	AddWriter(os.Stdout, WithFilter(Monitor))
	// which is equivalent to
	// 		AddWriter(os.Stdout, WithFilter(P | E | W))
}

func ExampleStringToType() {
	AddWriter(os.Stdout, WithFormatter(FormatWithColours), WithFilter(StringToType("PEW")))
	// which is equivalent to
	// 		AddWriter(os.Stdout, WithFilter(P | E | W))
}

func ExampleLabelToType() {
	// available labels:
	// 	none     // -
	// 	panic    // Panic
	// 	error    // Panic, Error
	// 	warning  // Panic, Error, Warning
	// 	info     // Panic, Error, Warning, Info
	// 	success  // Panic, Error, Warning, Info, Success
	// 	debug    // Panic, Error, Warning, Info, Success, and Debug
	// 	critical // Panic, Error
	// 	monitor  // Panic, Error, Warning
	// 	verbose  // Panic, Error, Warning, Info, Success
	// 	all      // Panic, Error, Warning, Info, Success, and Debug
	AddWriter(os.Stdout, WithFormatter(FormatWithColours), WithFilter(LabelToType("success")))
	// which is equivalent to
	// 		AddWriter(os.Stdout, WithFilter(P | E | W | I | S)) // debug logs are filtered out
}
