package round

import (
	"io"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

// Terminal escape sequences.
var (
	hideCursor     = []byte{27, '[', '?', '2', '5', 'l'}
	showCursor     = []byte{27, '[', '?', '2', '5', 'h'}
	saveCursor     = []byte{27, '[', 's'}
	clearFromSaved = []byte{27, '[', 'u', 27, '[', 'K'}
)

// FileWriter is an io.Writer that also has an Fd.
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// SpinMe goes right round. It's an io.WriteCloser.
type SpinMe struct {
	out  FileWriter
	now  string
	mu   *sync.Mutex
	tick *time.Ticker
}

// NewSpinMe creates a SpinMe and sets it spinning. It spins until it is closed.
// If the FileWriter is not a terminal, the spinner is bypassed.
func NewSpinMe(out FileWriter, s Style) SpinMe {
	if !terminal.IsTerminal(int(out.Fd())) || len(s.Frames) == 0 || s.Rate == time.Duration(0) {
		return SpinMe{out, "", nil, nil}
	}
	u := SpinMe{out, s.Frames[0], &sync.Mutex{}, time.NewTicker(s.Rate)}
	u.out.Write(append(append(saveCursor, hideCursor...), u.now...))
	go u.writeRound(s.Frames)
	return u
}

// writeRound spins the spinner right round. Like a record, baby.
func (u *SpinMe) writeRound(baby []string) {
	var f int
	for _ = range u.tick.C {
		f = (f + 1) % len(baby)
		u.now = baby[f]
		u.mu.Lock()
		u.out.Write(append(clearFromSaved, u.now...))
		u.mu.Unlock()
	}
}

// Write writes to out, moving the spinner to the end of what's written.
func (u *SpinMe) Write(p []byte) (int, error) {
	u.mu.Lock()
	u.out.Write(clearFromSaved)
	n, err := u.out.Write(p)
	u.out.Write(append(saveCursor, u.now...))
	u.mu.Unlock()
	return n, err
}

// Close will stop and remove the spinner.
func (u *SpinMe) Close() error {
	u.tick.Stop()
	u.out.Write(append(clearFromSaved, showCursor...))
	return nil
}