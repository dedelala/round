package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dedelala/round"
)

func main() {
	s := round.Pipe
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "block":
			s = round.Block
		case "cylon":
			s = round.Cylon
		case "hearts":
			s = round.Hearts
		case "moon":
			s = round.Moon
		case "help":
			fmt.Fprintf(os.Stderr, "Usage: %v [style]\n", os.Args[0])
		}
	}

	round.Go(s)
	io.Copy(round.Stdout, os.Stdin)
	round.Stop()
}
