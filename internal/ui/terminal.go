package ui

import (
	"os"

	"golang.org/x/term"
)

var openTTY = func() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0)
}

// IsInteractiveTerminal reports whether the process has a usable terminal for
// interactive prompts and full-screen UI.
func IsInteractiveTerminal() bool {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) || !term.IsTerminal(int(os.Stderr.Fd())) {
		return false
	}

	tty, err := openTTY()
	if err != nil {
		return false
	}
	return tty.Close() == nil
}
