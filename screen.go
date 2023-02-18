package screen

import (
	"bufio"
	"errors"
	"fmt"
	"golang.org/x/term"
	"os"
)

type Screen struct {
	orig          *term.State
	input         []byte
	output        *bufio.Writer
	Width, Height int
}

const (
	STYLE_NONE    = 0
	STYLE_BOLD    = 1
	STYLE_REVERSE = 7

	COLOR_RED     = 31
	COLOR_BLUE    = 34
	COLOR_MAGENTA = 35
	COLOR_CYAN    = 36
)

func New() (Screen, error) {
	var err error
	var screen Screen

	if !term.IsTerminal(0) {
		return screen, errors.New("this is not a terminal")
	}

	screen.orig, err = term.MakeRaw(0)
	if err != nil {
		return screen, err
	}

	screen.Width, screen.Height, err = term.GetSize(0)
	if err != nil {
		return screen, err
	}

	screen.input = make([]byte, 1)
	screen.output = bufio.NewWriter(os.Stdout)
	return screen, err
}

func (screen Screen) Write(p []byte) (int, error) {
	return screen.output.Write(p)
}

func (screen *Screen) Flush() error {
	return screen.output.Flush()
}

func (screen *Screen) HideCursor() {
	fmt.Fprint(screen, "\x1b[?25l")
}

func (screen *Screen) ShowCursor() {
	fmt.Fprint(screen, "\x1b[?25h")
}

func (screen *Screen) MoveCursor(x, y int) {
	fmt.Fprintf(screen, "\x1b[%d;%dH", y+1, x+1)
}

func (screen *Screen) Clear() {
	fmt.Fprintf(screen, "\x1b[2J\x1b[H\x1b[3J")
}

func (screen *Screen) Apply(effects ...int) {
	fmt.Fprint(screen, "\x1b[")
	for i, effect := range effects {
		if i > 0 {
			fmt.Fprint(screen, ";")
		}

		fmt.Fprint(screen, effect)
	}
	fmt.Fprint(screen, "m")
}

func (screen *Screen) Reset() {
	term.Restore(0, screen.orig)
	screen.Clear()
	screen.ShowCursor()
	screen.Apply(STYLE_NONE)
	screen.Flush()
}

func (screen *Screen) Input() (byte, error) {
	n, err := os.Stdin.Read(screen.input)
	if err != nil {
		return byte(0), err
	}

	if n == 0 {
		return byte(0), errors.New("could not read byte")
	}

	return screen.input[0], nil
}
