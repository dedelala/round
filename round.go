// Package round is a command line spinner. Start one with Go.
//
// The package will intelligently decide whether to write spinners on stdout,
// stderr or neither, depending if a terminal is present.
//
// Wrappers for Stdout and Stderr are provided so as not to interfere with the
// spinner while running.
package round

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	// Stdout is a spin-safe version of os.Stdout.
	Stdout io.Writer

	// Stderr is a spin-safe version of os.Stderr.
	Stderr io.Writer

	// Spinner frames and control bytes will be written on out.
	out *writer

	// spin is the current spinner.
	spin *spinMe

	// mu ensures clean output.
	mu = &sync.Mutex{}

	// terminal control bytes.
	hide  = []byte{27, '[', '?', '2', '5', 'l'}
	show  = []byte{27, '[', '?', '2', '5', 'h'}
	save  = []byte{27, '7'}
	load  = []byte{27, '8'}
	clear = []byte{27, '[', 'K'}
)

// Go makes a spinner go. It will Stop first if there is one running already.
func Go(s Style) {
	if out == nil {
		return
	}
	if spin != nil {
		Stop()
	}
	spin = &spinMe{s.Frames[0], make(chan bool)}
	out.writeAll(save, hide)
	go spin.writeRound(s.Frames, time.NewTicker(s.Rate))
}

// Stop will stop and remove the spinner.
func Stop() {
	if out == nil || spin == nil {
		return
	}
	spin.stop <- true
	out.writeAll(clear, show)
}

// writer is for writing out control bytes
type writer struct {
	out io.Writer
}

// write all writes all the bytes, or freaks the heck out
func (w *writer) writeAll(b ...[]byte) {
	for _, v := range b {
		if _, err := w.out.Write(v); err != nil {
			panic(err)
		}
	}
}

// blockingWriter will block on spin's writeRound
type blockingWriter struct {
	out io.Writer
}

// Write writes out, moving the spinner to the end of what's written.
func (w *blockingWriter) Write(p []byte) (int, error) {
	mu.Lock()
	out.writeAll(clear)
	n, err := w.out.Write(p)
	out.writeAll(save, []byte(spin.now), load)
	mu.Unlock()
	return n, err
}

// spinMe goes right round.
type spinMe struct {
	now  string
	stop chan bool
}

// writeRound spins the spinner right round. Like a record, baby.
func (u *spinMe) writeRound(baby []string, rightRound *time.Ticker) {
	var f int
	for {
		select {
		case <-rightRound.C:
			f = (f + 1) % len(baby)
			mu.Lock()
			u.now = baby[f]
			out.writeAll(clear, save, []byte(u.now), load)
			mu.Unlock()
		case <-u.stop:
			rightRound.Stop()
		}
	}
}

// init sets up globals and a goroutine to catch interrupts.
func init() {
	o := terminal.IsTerminal(int(os.Stdout.Fd()))
	e := terminal.IsTerminal(int(os.Stderr.Fd()))
	switch {
	case o && e:
		Stdout = &blockingWriter{os.Stdout}
		Stderr = &blockingWriter{os.Stderr}
		out = &writer{os.Stdout}
	case o:
		Stdout = &blockingWriter{os.Stdout}
		Stderr = os.Stderr
		out = &writer{os.Stdout}
	case e:
		Stdout = os.Stdout
		Stderr = &blockingWriter{os.Stderr}
		out = &writer{os.Stderr}
	default:
		Stdout = os.Stdout
		Stderr = os.Stderr
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		Stop()
	}()
}
