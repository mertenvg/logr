package logr

const (
	// ColourReset code
	ColourReset = "\x1B[0m"
)

type MetaData Meta

// Message used to send log message to logger goroutine
type Message struct {
	Type Type          `json:"type"`
	Time string        `json:"time"`
	Code string        `json:"code"`
	Desc string        `json:"description"`
	Meta MetaData      `json:"metadata,omitempty"`
	done chan struct{} `json:"-"`
}

// Reset the message object for later reuse
func (m *Message) Reset() {
	m.Type = None
	m.Time = ""
	m.Code = ""
	m.Desc = ""
	m.Meta = nil
	m.done = nil
}
