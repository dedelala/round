package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dedelala/round"
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
}

func main() {
	round.Go(style())

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		round.Stop()
		os.Exit(130)
	}()

	_, err := io.Copy(os.Stdout, os.Stdin)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
	}
	round.Stop()
}

func style() round.Style {
	flag.Usage = usage
	w := flag.Int("w", 8, "field width of a scroller or bouncer")
	f := flag.String("f", "[%s]", "format for a scroller or bouncer frame")
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
		os.Exit(2)
	case "block":
		return round.Block
	case "circle":
		return round.Circle
	case "cylon":
		return round.Cylon
	case "hearts":
		return round.Hearts
	case "moon":
		return round.Moon
	case "pipe":
		return round.Pipe
	case "wave":
		return round.Wave
	case "bounce":
		return round.NewBounce(*w, *f, msg)
	case "scroll":
		return round.NewRtL(*w, *f, msg)
	case "rtl":
		return round.NewRtL(*w, *f, msg)
	case "ltr":
		return round.NewLtR(*w, *f, msg)
	}

	return round.Pipe
}
