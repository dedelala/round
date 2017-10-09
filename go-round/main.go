package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dedelala/round"
	exit "github.com/dedelala/sysexits"
)

func usage() {
	f := `%v - copies stdin to stdout and shows a spinner
Usage: %[1]v [style]

The default style is pipe.

Style   | =  | Unicode Set
--------|----|--------------
block   | █  | 2580—259F Block Elements
cylon   | @  | 0020—007F Basic Latin
hearts  | 💖 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
moon    | 🌓 | 1F300—1F5FF Miscellaneous Symbols and Pictographs
pipe    | -  | 0020—007F Basic Latin

Scrollers
Usage: %[1]v [options] [scroll|bounce] [message...]

`
	fmt.Fprintf(os.Stderr, f, os.Args[0])
	flag.PrintDefaults()
	os.Exit(exit.Usage)
}

func main() {
	round.Go(divineStyle())
	io.Copy(round.Stdout, os.Stdin)
	round.Stop()
}

func divineStyle() round.Style {
	flag.Usage = usage
	w := flag.Int("w", 8, "field width of a scroller or bouncer")
	f := flag.String("f", "[%v]", "format for a scroller or bouncer frame")
	flag.Parse()

	if flag.NArg() == 0 {
		return round.Pipe
	}

	msg := "ooo"
	if flag.NArg() > 1 {
		msg = strings.Join(flag.Args()[1:], " ")
	}

	switch flag.Arg(0) {
	case "help":
		flag.Usage()
	case "block":
		return round.Block
	case "cylon":
		return round.Cylon
	case "hearts":
		return round.Hearts
	case "moon":
		return round.Moon
	case "pipe":
		return round.Pipe
	case "bounce":
		if *w < 0 {
			return round.NewInvertedBouncer(-*w, *f, msg)
		}
		return round.NewBouncer(*w, *f, msg)
	case "scroll":
		if *w < 0 {
			return round.NewInvertedScroller(-*w, *f, msg)
		}
		return round.NewScroller(*w, *f, msg)
	}

	return round.Pipe
}
