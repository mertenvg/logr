package logr

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type writerMessage struct {
	c *WriterConfig
	w io.Writer
}

var (
	mutex        = sync.RWMutex{}
	addWriter    = make(chan writerMessage)
	removeWriter = make(chan writerMessage)
)

// Formatter is a function that returns a formatted byte array representation of the given
// Message.
type Formatter func(m *Message) []byte

// FormatJSON is a Formatter that converts a Message to json
func FormatJSON(m *Message) []byte {
	r, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("failed to marshal Message in logr.FormatJSON: %v\n", err)
	}
	return append(r, []byte("\n")...)
}

// LogrFormat is a Formatter that converts a Message to a default format.
func FormatDefault(m *Message) []byte {
	var s string
	if m.Meta == nil {
		s = fmt.Sprintf(
			"%-25s | %s | %s | %s\n",
			m.Time,
			m.Code,
			m.Type.Rune(),
			m.Desc,
		)
	} else {
		s = fmt.Sprintf(
			"%-25s | %s | %s | %s | %+v\n",
			m.Time,
			m.Code,
			m.Type.Rune(),
			m.Desc,
			m.Meta,
		)
	}
	return []byte(s)
}

// FormatWithColours is a Formatter that converts a Message to a colourful default format.
func FormatWithColours(m *Message) []byte {
	var s string
	if m.Meta == nil {
		s = fmt.Sprintf(
			m.Type.Colour()+"%-25s | %s | %s | "+ColourReset+"%s\n",
			m.Time,
			m.Code,
			m.Type.Rune(),
			m.Desc,
		)
	} else {
		s = fmt.Sprintf(
			m.Type.Colour()+"%-25s | %s | %s | "+ColourReset+"%s"+m.Type.Colour()+" | %+v"+ColourReset+"\n",
			m.Time,
			m.Code,
			m.Type.Rune(),
			m.Desc,
			m.Meta,
		)
	}

	return []byte(s)
}

type WriterConfig struct {
	format Formatter
	filter Type
}

type WriterConfigModifier func(c WriterConfig) WriterConfig

// WithFormatter creates a WriterConfigModifier that defines how a log message is formatted on the configuration
// for a Writer. A custom formatter may be provided to convert the Message to any desired format before it is passed
// to the Writer.
func WithFormatter(f Formatter) WriterConfigModifier {
	return func(oc WriterConfig) WriterConfig {
		oc.format = f
		return oc
	}
}

// WithFilter creates a WriterConfigModifier that sets a filter on the configuration for a Writer
func WithFilter(f Type) WriterConfigModifier {
	return func(oc WriterConfig) WriterConfig {
		oc.filter = f
		return oc
	}
}

// AddWriter add a io.Writer to the collection of writers that store the log messages.
func AddWriter(w io.Writer, configs ...WriterConfigModifier) (stop func()) {
	// default config
	oc := WriterConfig{
		format: FormatDefault,
		filter: All,
	}
	// apply optional extra config modifiers
	for _, c := range configs {
		oc = c(oc)
	}
	// register the output
	wm := writerMessage{
		w: w,
		c: &oc,
	}
	addWriter <- wm
	return func() {
		removeWriter <- wm
	}
}
