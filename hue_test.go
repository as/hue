package hue

import (
	"testing"
	"fmt"
)

func TestEncodeDecode(t *testing.T) {
	var h Hue

	testInputs := []string{
		"The quick brown fox jumped over the lazy dog", // All letters of the alphabet
		"Test",               // Less than 10 characters, but not empty
		"",                   // Empty
		"\300\r\n\t\x0d\x0a", // Non-printables
		"\033[00;00m\033[0m", // Characters that mimic the ECMA encoding
	}

	/*
	 * Test the input strings with every possible foreground / background color
	 * by encoding and subsequently decoding the and comparing
	 * the result to the original test input.
	 */
	for _, u := range testInputs {
		for _, v := range Names {
			for _, w := range Names {
				hs := Encode(&h, u)
				if u != hs.Decode() {
					t.Log(fmt.Sprintf("%s != %s", u, hs.Decode()))
					t.Fail()
				}
				h.Fg = w
				h.Bg = v
			}
		}
	}
}
