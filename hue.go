// Package hue implements terminal background and foreground output colorization.
// A Hue object provides a wrapper for output functions (such as fmt.Printf). A Writer
// object allows an io.Writer to bind to a Hue object for colored writes.
package hue

import (
	"fmt"
	"io"
)

// Foreground color codes
const (
	Black = iota + 30
	Red
	Green
	Brown
	Blue
	Magenta
	Cyan
	White
	Default
)

// Names is a map of Hue strings to color codes
var Names = map[string]int{
	"Black":   Black,
	"Red":     Red,
	"Green":   Green,
	"Brown":   Brown,
	"Blue":    Blue,
	"Magenta": Magenta,
	"Cyan":    Cyan,
	"White":   White,
	"Default": Default,
}

// SetHue sets the Writer's hue
func (w *Writer) SetHue(h *Hue) {
	w.Hue = h
}

// NewWriter returns a new Writer with the Hue 'h'
func NewWriter(w io.Writer, h *Hue) *Writer {
	n := new(Writer)
	n.wrapped = w
	n.SetHue(h)
	return n
}

// Write colorizes and writes the contents of p to the underlying
// writer object.
func (w Writer) Write(p []byte) (n int, err error) {
	return w.wrapped.Write([]byte(w.Sprintf("%s", string(p))))
}

// WriteString colorizes and writes the string s to the
// underlying writer object
func (w Writer) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}

// Writer implements colorization for an underlying io.Writer object
type Writer struct {
	*Hue
	wrapped io.Writer
}

// String is a string containing ECMA-48 color codes. Its purpose is to
// remind the user at compile time that it differs from the string builtin.
type String string

// Hue holds the foreground color and background color as integers
type Hue struct {
	Fg, Bg int
}

// isUnset returns true if the hue integer value is 0
func isUnset(i int) bool {
	return i == 0
}

// Decode strips all color data from the String object
// and returns a standard string
func (hs String) Decode() (s string) {
	if l := len(hs); l < 12 {
		panic(fmt.Sprintf("Can't decode hue.String: \"%s\" because it's length is \"%d\" (minimum length is 12)", hs, l ))
	}

	b := []byte(hs)
	return string(b[8:len(b) - 4])
}

// Encode encapsulates interface a's string representation
// with the ECMA-40 color codes stored in the Hue structure.
func Encode(h *Hue, a interface{}) String {
	if isUnset(h.Fg) {
		h.Fg = Default
	}

	if isUnset(h.Bg) {
		h.Bg = Default
	}

	return String(fmt.Sprintf("\033[%d;%dm%v\033[0m", h.Fg, h.Bg+10, a))
}

// Sprintf behaves like fmt.Sprintf, except it colorizes the output String
func (h *Hue) Sprintf(format string, a interface{}) String {
	return String(fmt.Sprintf(string(Encode(h, format)), a))
}

// Printf behaves like fmt.Printf, except it colorizes the output
func (h *Hue) Printf(format string, a interface{}) {
	fmt.Printf(string(Encode(h, format)), a)
}

// Print behaves like fmt.Print, except it colorizes the output
func (h *Hue) Print(a interface{}) {
	fmt.Print(Encode(h, a))
}

// Println behaves like fmt.Println, except it colorizes the output
func (h *Hue) Println(a interface{}) {
	fmt.Println(Encode(h, a))
}
