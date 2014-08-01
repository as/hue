// Package hue implements terminal background and foreground output colorization.
// A hue object provides a wrapper for output functions (such as fmt.Printf). A Writer
// object allows an io.Writer to bind to a hue object for colored writes. A RegexpWriter
// implements automatic colorization via regular expressions.
package hue

import (
	"fmt"
	"io"
	"regexp"
)

// Foreground color codes
const (
	Black = iota + 30
	Red
	Green
	Brown // Looks more like yellow on most terminal emulators
	Blue
	Magenta
	Cyan
	White
	Underline
	Default
)

const (
	// ASCIIFmt is a format specifer for a ECMA-48 color string
	ASCIIFmt = "\033[%d;%dm"
	// ASCIIReset is a reset code that restores the colors to their defaults
	ASCIIReset = "\033[0m"
	// ASCIIFmtReset is a combination of ASCIIFmt and ASCIIReset with a value is sandwiched in between.
	ASCIIFmtReset = "\033[%d;%dm%v\033[0m"
)

// Names is a map of hue strings to color codes
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
func (w *Writer) SetHue(h *hue) {
	w.hue = h
}

// NewWriter returns a new Writer with the hue 'h'
func NewWriter(w io.Writer, h *hue) *Writer {
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
	*hue
	wrapped io.Writer
}

// String is a string containing ECMA-48 color codes. Its purpose is to
// remind the user at compile time that it differs from the string builtin.
type String string

// New creates a new hue object with foreground and background colors specified.
func New(fg, bg int) *hue {
	h := new(hue)
	h.SetFg(fg)
	h.SetBg(bg)
	return h
}

func (h *hue) SetFg(c int) {
	h.fg = c
}

func (h *hue) SetBg(c int) {
	h.bg = c + 10
}

func (h *hue) Fg() int {
	return h.fg
}

func (h *hue) Bg() int {
	return h.bg
}

// hue holds the foreground color and background color as integers
type hue struct {
	fg, bg int
}

// Decode strips all color data from the String object
// and returns a standard string
func (hs String) Decode() (s string) {
	if l := len(hs); l < 12 {
		panic(fmt.Sprintf("Can't decode hue.String: \"%s\" because it's length is \"%d\" (minimum length is 12)", hs, l))
	}

	b := []byte(hs)
	return string(b[8 : len(b)-4])
}

// Encode encapsulates interface a's string representation
// with the ECMA-40 color codes stored in the hue structure.
func Encode(h *hue, a interface{}) String {
	return String(fmt.Sprintf(ASCIIFmtReset, h.Fg(), h.Bg(), a))
}
// Sprintf behaves like fmt.Sprintf, except it colorizes the output String
func (h *hue) Sprintf(format string, a interface{}) String {
	return String(fmt.Sprintf(string(Encode(h, format)), a))
}

// Printf behaves like fmt.Printf, except it colorizes the output
func (h *hue) Printf(format string, a interface{}) {
	fmt.Printf(string(Encode(h, format)), a)
}

// Print behaves like fmt.Print, except it colorizes the output
func (h *hue) Print(a interface{}) {
	fmt.Print(Encode(h, a))
}

// Println behaves like fmt.Println, except it colorizes the output
func (h *hue) Println(a interface{}) {
	fmt.Println(Encode(h, a))
}

// RegexpWriter implements colorization for a io.Writer object by processing
// a set of rules. Rules are hue objects assocated with regular expressions.
type RegexpWriter struct {
	rules   []rule
	wrapped io.Writer
}

type rule struct {
	*hue
	*regexp.Regexp
}

// NewRegexpWriter returns a new RegexpWriter
func NewRegexpWriter(w io.Writer) *RegexpWriter {
	n := new(RegexpWriter)
	n.wrapped = w
	return n
}

// AddRuleStringPOSIX binds a hue to the POSIX regexp in the string 's'.
// Similar to AddRule, except the caller passes in an uncompiled POSIX regexp.
func (w *RegexpWriter) AddRuleStringPOSIX(h *hue, s string) {
	re := regexp.MustCompilePOSIX(s)
	w.AddRule(h, re)
}

// AddRuleString binds a hue to the regexp in the string 's'.
// Similar to AddRule, except the caller passes in an uncompiled regexp.
func (w *RegexpWriter) AddRuleString(h *hue, s string) {
	re := regexp.MustCompile(s)
	w.AddRule(h, re)
}

// AddRule binds a hue to a regular expression.
func (w *RegexpWriter) AddRule(h *hue, re *regexp.Regexp) {
	//w.rules.PushBack( rule{h, re} )
	w.rules = append(w.rules, rule{h, re})
}

// FlushRules deletes all rules added with AddRule from Writer
func (w *RegexpWriter) FlushRules() {
	w.rules = nil
}

// PrintRules prints out the rules
func (w *RegexpWriter) PrintRules() {
	for _, v := range w.rules {
		fmt.Println(v)
	}
}

// WriteString is similar to Write, except it writes a string to the underlying
// buffer instead of a byte slice.
func (w RegexpWriter) WriteString(s string) (n int, err error) {
	return w.Write([]byte(s))
}

// Write writes the contents of p into the buffer after processesing the regexp
// rules added to Writer with AddRule. Write colorizes the contents as it writes
// to the underlying writer object.
func (w RegexpWriter) Write(p []byte) (n int, err error) {
	huemap := make([]byte, len(p))
	rulemap := make([]*hue, len(w.rules)+1)
	rulemap[0] = &hue{}

	for i := 1; i < len(rulemap); i++ {
		r := w.rules[i-1]
		x := r.FindAllIndex(p, -1)

		rulemap[i] = r.hue
		for _, w := range x {
			for j := w[0]; j < w[1]; j++ {
				huemap[j] = byte(i)
			}
		}
	}

	var hue byte
	for i := range p {
		if huemap[i] != hue {
			hue = huemap[i]
			th := rulemap[hue]

			nb, err := fmt.Fprintf(w.wrapped, ASCIIFmt, th.Fg(), th.Bg())
			if err != nil {
				return n, err
			}
			n += nb
		}

		nb, err := fmt.Fprintf(w.wrapped, "%c", p[i])
		if err != nil {
			return n, err
		}
		n += nb
	}
	fmt.Print(ASCIIReset)

	return n, err
}
