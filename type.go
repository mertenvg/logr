package logr

import (
	"fmt"
	"strings"
)

// Type of log item
type Type int

// Available log types
const (
	None Type = iota      // None
	P    Type = 1 << iota // Panic
	E                     // Error
	W                     // Warning
	I                     // Info
	D                     // Debug
	S                     // Success

	Critical = P | E
	Monitor  = Critical | W
	Verbose  = Monitor | I | S
	All      = Verbose | D
)

var (
	typeRuneMap = map[Type]string{
		None: "-",
		P:    "P",
		E:    "E",
		W:    "W",
		I:    "I",
		D:    "D",
		S:    "S",
	}
	typeStringMap = map[Type]string{
		None: "none",
		P:    "panic",
		E:    "error",
		W:    "warning",
		I:    "info",
		D:    "debug",
		S:    "success",
	}
	typeColourMap = map[Type]string{
		None: "\x1B[0m",
		P:    "\x1B[38;5;124m",
		E:    "\x1B[38;5;124m",
		W:    "\x1B[38;5;208m",
		I:    "\x1B[38;5;33m",
		D:    "\x1B[38;5;153m",
		S:    "\x1B[38;5;34m",
	}
	typeLabelMap = map[string]Type{
		"none":     None,
		"panic":    P,
		"error":    Critical,
		"warning":  Monitor,
		"info":     Monitor | I,
		"success":  Verbose,
		"debug":    All,
		"critical": Critical,
		"monitor":  Monitor,
		"verbose":  Verbose,
		"all":      All,
	}
)

// String returns a descriptive string of the log Type
func (t Type) String() string {
	if s, ok := typeStringMap[t]; ok {
		return s
	}
	return "unknown"
}

// MashalJSON implements json.Marshaller
func (t Type) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", t.String())), nil
}

// Rune returns the rune of the log Type
func (t Type) Rune() string {
	if s, ok := typeRuneMap[t]; ok {
		return s
	}
	return fmt.Sprintf("Type(%d)", t)
}

// Colour returns the bash colour code for the log Type
func (t Type) Colour() string {
	if s, ok := typeColourMap[t]; ok {
		return s
	}
	return typeColourMap[None]
}

// IntToType converts and integer to a log type/level
func IntToType(i int) Type {
	return Type(i)
}

// StringToType converts a string of codes to a log type/level
func StringToType(s string) Type {
	t := Type(0)
	for _, r := range s {
		t |= RuneToType(r)
	}
	return t
}

var (
	types = "PEWIDS"
)

// RuneToType converts a code to a log type/level
func RuneToType(r rune) Type {
	i := strings.IndexRune(types, r)
	if i == -1 {
		return IntToType(0)
	}
	return IntToType(i)
}

// LabelToType converts a label to a log type/level
// Available labels are one of: none, panic, error, warning, info, success, debug, critical, monitor, verbose, all
func LabelToType(l string) Type {
	t, ok := typeLabelMap[strings.TrimSpace(strings.ToLower(l))]
	if !ok {
		panic(fmt.Sprintf("logr label `%s` not found in supported types (none, panic, error, warning, info, success, debug, critical, monitor, verbose, all)", l))
	}
	return t
}
