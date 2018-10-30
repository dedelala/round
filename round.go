// Package round is a command line spinner. Start one with Go.
//
// The package will decide whether to write the spinner to stdout, stderr or
// neither, depending if a terminal is present.
//
// Wrappers for Stdout and Stderr are provided to prevent interference with the
// spinner while running.
package round

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/pkg/term"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	hide  = "\x1b[?25l"
	show  = "\x1b[?25h"
	save  = "\x1b7"
	load  = "\x1b8"
	clear = "\x1b[K"
	erase = "\x1b[J"
)

var (
	bars     = make(chan bar)
	enabled  = terminal.IsTerminal(int(os.Stderr.Fd()))
	progress func(bar)
	stop     func()
)

type bar struct {
	k, v string
}

func update(bs []bar, b bar) []bar {
	if b.k == "" && b.v == "" {
		return bs
	}
	for i, v := range bs {
		if v.k == b.k {
			bs[i] = b
			return bs
		}
	}
	return append(bs, b)
}

func pass(recv <-chan bar, send chan<- bar) {
	bs := []bar{}
	for {
		var b bar
		if len(bs) > 0 {
			b = bs[0]
		}
		select {
		case b := <-recv:
			bs = update(bs, b)
		case send <- b:
			if len(bs) > 0 {
				bs = bs[1:]
			}
		}
	}
}

func round(s Style, bars chan bar, stop, done chan struct{}) {
	var (
		frame string
		clk   <-chan time.Time
		sig   = make(chan os.Signal)
		bs    = []bar{}
	)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	os.Stderr.WriteString(hide)

	draw := func() {
		d := erase + frame
		for _, b := range bs {
			d += "\n" + b.k + " " + b.v
		}
		d += "\n\x1b[" + strconv.Itoa(len(bs)+1) + "A"
		os.Stderr.WriteString(d)
	}

	for {
		if clk == nil {
			f, d := s.Frame()
			frame = f
			draw()
			clk = time.After(d)
		}
		select {
		case <-clk:
			clk = nil
		case b := <-bars:
			bs = update(bs, b)
		case <-stop:
			os.Stderr.WriteString(erase + show)
			close(done)
			return
		case <-sig:
			os.Stderr.WriteString(erase + show)
			close(done)
			return
		}
	}
}

// Go starts a spinner on stderr with the given style. Any running spinner will
// be stopped.
func Go(s Style) {
	if !enabled {
		return
	}
	Stop()
	sc, dc := make(chan struct{}), make(chan struct{})
	stop = func() {
		select {
		case sc <- struct{}{}:
		case <-dc:
		}
		<-dc
	}

	go round(s, bars, sc, dc)
}

func Progress(k, v string) {
	if progress == nil {
		recv := make(chan bar)
		progress = func(b bar) {
			recv <- b
		}
		go pass(recv, bars)
	}
	progress(bar{k, v})
}

// Stop stops any running spinner.
func Stop() {
	if stop == nil {
		return
	}
	stop()
}

//type Reader struct {
//io.Reader
//Label string
//}

//func (r *Reader) Read(p []byte) (int, error) {
//n, err := r.Reader.Read(p)
//}

func getpos(t *term.Term) (int, int, error) {
	if err := t.SetCbreak(); err != nil {
		return -1, -1, err
	}
	t.Write([]byte("\x1b[6n"))
	var r, c int
	n, err := fmt.Fscanf(t, "\x1b[%d;%dR", &r, &c)
	if err != nil {
		return -1, -1, err
	}
	if n != 2 {
		return -1, -1, fmt.Errorf("bad scan n=%d", n)
	}
	return r, c
}

func ttyname() (string, error) {
	info, err := os.Stderr.Stat()
	if err != nil {
		return "not a tty", fmt.Errorf("stat stderr: %s", err)
	}
	sys, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "not a tty", fmt.Errorf("stat stderr: not Stat_t (%T)", info.Sys())
	}
	dev := sys.Rdev
	match := ""

	wf := func(path string, info os.FileInfo, err error) error {
		if match != "" {
			return filepath.SkipDir
		}

		rel, _ := filepath.Rel("/dev", path)
		dir, _ := filepath.Split(rel)
		if !(dir == "" || dir == "/pts") {
			return filepath.SkipDir
		}

		sys, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return nil
		}
		if sys.Rdev != dev {
			return nil
		}

		match = path
		return filepath.SkipDir
	}

	if err := filepath.Walk("/dev", wf); err != nil {
		return "not a tty", fmt.Errorf("walk /dev: %s", err)
	}

	if match == "" {
		return "not a tty", nil
	}
	return match, nil
}
