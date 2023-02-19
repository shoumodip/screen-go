# Screen
Simple Terminal UI Library in Go

## Installation
```console
$ go get github.com/shoumodip/screen-go
```

## Usage
```go
package main

import (
	"fmt"
	"github.com/shoumodip/screen-go"
)

func main() {
	scr, err := screen.New()
	if err != nil {
		panic(err)
	}
	defer scr.Reset()

	scr.HideCursor()
	defer scr.ShowCursor()

	for {
		scr.Clear()

		// Coordinates start at (0, 0) instead of (1, 1)
		scr.MoveCursor(0, 0)
		scr.Apply(screen.STYLE_BOLD, screen.COLOR_BLUE)

		// Newlines have to be '\r\n'
		fmt.Fprint(scr, "Hello, world!\r\n")
		scr.Apply(screen.STYLE_NONE)

		// Present the rendered state to the screen
		scr.Flush()

		key, err := scr.Input()
		if err != nil {
			panic(err)
		}

		if key == 'q' {
			break
		}
	}
}
```
