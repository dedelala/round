package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dedelala/round"
	exit "github.com/dedelala/sysexits"
)

const usage = `%v - copies stdin to stdout and shows a spinner
Usage: %[1]v [style]

The default style is pipe.

Style   | =  | Unicode Set
--------|----|--------------
block   | █  | 2580—259F Block Elements
cylon   | @  | 0020—007F Basic Latin
hearts  | 💖 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
moon    | 🌓 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
pipe    | -  | 0020—007F Basic Latin
`

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
			fmt.Fprintf(os.Stderr, usage, os.Args[0])
			os.Exit(exit.Usage)
		}
	}

	round.Go(s)
	io.Copy(round.Stdout, os.Stdin)
	round.Stop()
}
