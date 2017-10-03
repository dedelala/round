// Package round is a command line spinner. Start one with Start.
package round

import (
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	term "golang.org/x/crypto/ssh/terminal"
)

var (
	// Stdout is a spin-safe version of os.Stdout
	Stdout io.Writer

	// Stderr is a spin-safe version of os.Stderr
	Stderr io.Writer

	// Spin Control
	out  io.Writer
	spin *spinMe
	mu   = spinit()

	// Terminal escape sequences.
	hide      = []byte{27, '[', '?', '2', '5', 'l'}
	show      = []byte{27, '[', '?', '2', '5', 'h'}
	save      = []byte{27, '[', 's'}
	clear     = []byte{27, '[', 'u', 27, '[', 'K'}
	saveHide  = append(save, hide...)
	clearShow = append(clear, show...)
)

// spinit is the init for the spin it.
func spinit() *sync.Mutex {
	o := term.IsTerminal(int(os.Stdout.Fd()))
	e := term.IsTerminal(int(os.Stderr.Fd()))
	switch {
	case o && e:
		Stdout = &blockingWriter{os.Stdout}
		Stderr = &blockingWriter{os.Stderr}
		out = os.Stdout
	case o:
		Stdout = &blockingWriter{os.Stdout}
		Stderr = os.Stderr
		out = os.Stdout
	case e:
		Stdout = os.Stdout
		Stderr = &blockingWriter{os.Stderr}
		out = os.Stderr
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

	return &sync.Mutex{}
}

// Start starts a spinner.
func Start(s Style) {
	if out == nil {
		return
	}
	if spin != nil {
		Stop()
	}
	spin = &spinMe{s.Frames[0], make(chan bool)}
	out.Write(saveHide)
	go spin.writeRound(s.Frames, time.NewTicker(s.Rate))
}

// Stop will stop and remove the spinner.
func Stop() {
	if out == nil || spin == nil {
		return
	}
	spin.stop <- true
	out.Write(clearShow)
}

type blockingWriter struct {
	out io.Writer
}

// Write writes to the underlying Writer, moving the spinner to the end of what's written.
func (w *blockingWriter) Write(p []byte) (int, error) {
	mu.Lock()
	w.out.Write(clear)
	n, err := w.out.Write(p)
	w.out.Write(append(save, spin.now...))
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
			out.Write(append(clear, u.now...))
			mu.Unlock()
		case <-u.stop:
			rightRound.Stop()
			break
		}
	}
}
