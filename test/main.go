package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/dedelala/round"
	"github.com/pkg/term"
)

func main() {
	log.SetFlags(0)

	t, err := term.Open(ttyname())
	if err != nil {
		log.Fatal(err)
	}

	r, c := getpos(t)
	log.Printf("r: %d\tc: %d", r, c)

	round.Go(round.Hearts)

	go func() {
		for i := 0; i < 10; i++ {
			r, c := getpos(t)
			log.Printf("r: %d\tc: %d", r, c)
			<-time.After(500 * time.Millisecond)
		}
	}()

	for _, s := range []string{"blep", "mlem", "nyaa"} {
		<-time.After(time.Second)
		round.Progress("cat", s)
		fmt.Println("PRINT ", s)
		log.Print(s)
	}
	<-time.After(time.Second)

	round.Stop()
}

func getpos(t *term.Term) (int, int) {
	if err := t.SetCbreak(); err != nil {
		log.Fatal(err)
	}
	t.Write([]byte("\x1b[6n"))
	var r, c int
	n, err := fmt.Fscanf(t, "\x1b[%d;%dR", &r, &c)
	if err != nil {
		log.Fatal(err)
	}
	if n != 2 {
		log.Fatal("scan failed n=", n)
	}
	return r, c
}

func ttyname() string {
	info, err := os.Stderr.Stat()
	if err != nil {
		return "not a tty: " + err.Error()
	}
	sys, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "no stat_t"
	}
	dev := sys.Rdev

	var match string

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
		if dev == sys.Rdev {
			match = path
			return filepath.SkipDir
		}
		return nil
	}

	filepath.Walk("/dev", wf)

	if match == "" {
		return "not a tty"
	}
	return match
}
