## Hue
Package hue implements terminal _colorization_ via ECMA-40 color codes. It provides four objects for working with colors: Hue, Writer, RegexpWriter and String.

## Quick tour
A Hue object provides a wrapper for output functions, such as fmt.Println.

```
	// Create a hue with a white foreground and red background
	red := hue.New(hue.White, hue.Red)

	// Print the red string with a hue object
	red.Println("Red")
```


A Writer object allows an io.Writer to bind to a Hue object for colored writes.

```
	green := hue.New(hue.Green, hue.Default)

	// Print a green string with a hue.Writer
	hw := hue.NewWriter(os.Stdout, green)
	hw.WriteString("Green\n")
```


A String object holds a colored string and can be created with Encode

```
	// Store a hue string as a variable
	var hs hue.String 
	hs = hue.Encode(&blue, "Always be dilligent and see your work through till th")

	fmt.Printf("Colored, even with fmt.Printf(): %s\n", hs)

	// A String object can be decoded to a builtin string using Decode
	fmt.Printf("Uncolored: %s\n", hs.Decode())
```


To automatically colorize output, create a hue.RegexpWriter
and add coloring rules with regular expressions.

```
	re := hue.NewRegexpWriter(os.Stdout)

	alarmState := `Weather control device online
    Snow storage silo: Warning Capacity Exceeded
    Snow storage silo: Danger: Pressure over limit`

	re.AddRuleStringPOSIX(green, ".*online.*")
	re.AddRuleStringPOSIX(yellow, "Warning")
	re.AddRuleStringPOSIX(red, ".*Danger.*")

	re.WriteString(alarmState)
```


### Full example
```
	package main
	
	import (
		"os"
	)
	
	import (
		"github.com/as/hue"
	)
	
	func main() {
		// Create a hue with a white foreground and red background
		red := hue.New(hue.White, hue.Red)
	
		// Create more hues
		green := hue.New(hue.Green, hue.Default)
		blue := hue.New(hue.Blue, hue.Default)
		yellow := hue.New(hue.Brown, hue.Default)
		magenta := hue.New(hue.Magenta, hue.Default)
	
		// Print the red string with a hue object
		red.Println("Red")
	
		// Print a green string with a hue.Writer
		hw := hue.NewWriter(os.Stdout, green)
		hw.WriteString("Green\n")
	
		// You can change the hue with SetHue
		hw.SetHue(blue)
		hw.WriteString("Blue\n")
	
		// To automatically colorize output, create a hue.RegexpWriter
		// and add coloring rules with regular expressions.
		re := hue.NewRegexpWriter(os.Stdout)
	
		alarmState := `
	Weather control device online
	Snow storage silo: Warning Capacity Exceeded
	Snow storage silo: Danger: Pressure over limit
	`
		re.AddRuleStringPOSIX(green, ".*online.*")
		re.AddRuleStringPOSIX(yellow, "Warning")
		re.AddRuleStringPOSIX(red, ".*Danger.*")
	
		re.WriteString(alarmState)
	
		// Advanced Example: This shows what you can do by layering regexp rules.
		// Delete all previously-added rules
		re.FlushRules()
	
		// Test file: /etc/passwd
		test := `
	bin:x:1:1:bin:/bin:/usr/bin/nologin
	daemon:x:2:2:daemon:/:/usr/bin/nologin
	mail:x:8:12:mail:/var/spool/mail:/usr/bin/nologin
	ftp:x:14:11:ftp:/srv/ftp:/usr/bin/nologin
	http:x:33:33:http:/srv/http:/usr/bin/nologin
	dbus:x:81:81:dbus:/:/usr/bin/nologin
	nobody:x:99:99:nobody:/:/usr/bin/nologin
	colord:x:124:124::/var/lib/colord:/bin/false
	dnsmasq:x:999:999:dnsmasq daemon:/:/sbin/nologin
	tiktoor::0:0::/:/bin/bash
	lxdm:x:121:121:Lightweight X11 Display Manager:/var/lib/lxdm:/sbin/nologin
	rpc:x:32:32:Rpcbind Daemon:/dev/null:/bin/false
	ntp:x:87:87:Network Time Protocol:/var/lib/ntp:/bin/false`
	
		// First, color every character magenta
		re.AddRuleStringPOSIX(magenta, ".")
	
		// Then, color every character up to the first colon blue
		re.AddRuleStringPOSIX(blue, "^[^:]+:")
	
		// Color every character from the final colon to the end of the line red
		re.AddRuleStringPOSIX(red, ":[^:]+*$")
	
		// If the text 'nologin' or 'false' is present in the last colon-seperated value, color it green
		re.AddRuleStringPOSIX(green, ":[^:]+(nologin|false)$")
	
		// Color all the colons yellow
		re.AddRuleStringPOSIX(yellow, ":")
	
		// Output the test string with according to the rules we just added.
		re.Write([]byte(test))
	}
	
```
