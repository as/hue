package main

import (
	"os"
	"fmt"
)

import (
	"github.com/as/hue"
)

func main() {
	// Step 1: Create a hue
	red := hue.Hue{
		Fg: hue.Red,
		Bg: hue.White,
	}

	// Step 2: Print a string in color
	red.Println("Red string is red /w white background")

	// Another way: Use a hue.Writer
	green := hue.Hue{hue.Green, hue.Default}
	blue := hue.Hue{hue.Blue, hue.Default}

	hw := hue.NewWriter(os.Stdout, &green)
	hw.Write([]byte("Write some green bytes to stdout\n"))

	// You can change the hue with SetHue
	hw.SetHue(&blue)
	hw.WriteString("Write this blue string to stdout\n")

	// Store a hue string as a variable
	var hs hue.String 
	hs = hue.Encode(&blue, "Always be dilligent and see your work through till th")

	// Print the hue.String with the "fmt" package
	fmt.Printf("Colored, even with fmt.Printf(): %s\n", hs)

	// Convert the hue.String back to a standard string
	fmt.Printf("Uncolored: %s\n", hs.Decode())
}
