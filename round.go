package round

import (
	"io"
	"sync"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

// Terminal escape sequences.
var (
	hide      = []byte{27, '[', '?', '2', '5', 'l'}
	show      = []byte{27, '[', '?', '2', '5', 'h'}
	save      = []byte{27, '[', 's'}
	clear     = []byte{27, '[', 'u', 27, '[', 'K'}
	saveHide  = append(save, hide...)
	clearShow = append(clear, show...)
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
	stop chan bool
}

// NewSpinMe creates a SpinMe and sets it spinning. It spins until it is closed.
// If the FileWriter is not a terminal, the spinner is bypassed.
func NewSpinMe(out FileWriter, s Style) SpinMe {
	if !terminal.IsTerminal(int(out.Fd())) || len(s.Frames) == 0 || s.Rate == time.Duration(0) {
		return SpinMe{out, "", nil, nil}
	}
	u := SpinMe{out, s.Frames[0], &sync.Mutex{}, make(chan bool)}
	u.out.Write(append(saveHide, u.now...))
	go u.writeRound(s.Frames, time.NewTicker(s.Rate))
	return u
}

// writeRound spins the spinner right round. Like a record, baby.
func (u *SpinMe) writeRound(baby []string, rightRound *time.Ticker) {
	var f int
	for {
		select {
		case <-rightRound.C:
			f = (f + 1) % len(baby)
			u.mu.Lock()
			u.now = baby[f]
			u.out.Write(append(clear, u.now...))
			u.mu.Unlock()
		case <-u.stop:
			rightRound.Stop()
			break
		}
	}
}

// Write writes to out, moving the spinner to the end of what's written.
func (u *SpinMe) Write(p []byte) (int, error) {
	if u.mu == nil {
		return u.out.Write(p)
	}
	u.mu.Lock()
	u.out.Write(clear)
	n, err := u.out.Write(p)
	u.out.Write(append(save, u.now...))
	u.mu.Unlock()
	return n, err
}

// Close will stop and remove the spinner.
func (u *SpinMe) Close() error {
	if u.mu == nil {
		return nil
	}
	close(u.stop)
	u.mu.Lock()
	u.out.Write(clearShow)
	u.mu.Unlock()
	return nil
}
