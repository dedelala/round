package round

import (
	"fmt"
	"time"
)

type Style interface {
	Frame() (string, time.Duration)
}

// Block is a block sequence from 2580—259F Block Elements.
var Block Style = &sequence{
	fs: []string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█", "▇", "▆", "▅", "▄", "▃", "▂", "▁", ""},
	d:  60 * time.Millisecond,
}

var Circle Style = &sequence{
	fs: []string{"◐", "◓", "◑", "◒"},
	d:  75 * time.Millisecond,
}

// Cylon is an ominous sequence using 0020—007F Basic Latin.
var Cylon = NewBounce(7, "\x1b[1m(\x1b[31m%v\x1b[0;1m)\x1b[0m", "@")

// Roll the Dice! Not a random style, an actual rolling dice.
// U+2600-U+26FF Miscellaneous Symbols.
var Wave Style = wave{}

// Hearts is a sequence of rainbow hearts. Clearly the best style!
// 1F300—1F5FF Misc Symbols and Pictographs.
var Hearts Style = &sequence{
	fs: []string{"💖💛💚💙💜", "💜💖💛💚💙", "💙💜💖💛💚", "💚💙💜💖💛", "💛💚💙💜💖"},
	d:  90 * time.Millisecond,
}

// Moon is a sequence of moon phases. It uses "🌕" from 1F300—1F5FF Misc Symbols and Pictographs.
var Moon Style = &sequence{
	fs: []string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
	d:  90 * time.Millisecond,
}

// Pipe is a sequence of pipes and slashes. The classic spinner.
var Pipe Style = &sequence{
	fs: []string{"|", "/", "-", "\\"},
	d:  60 * time.Millisecond,
}

type sequence struct {
	fs []string
	d  time.Duration
	i  int
}

func (s *sequence) Frame() (string, time.Duration) {
	s.i = (s.i + 1) % len(s.fs)
	return s.fs[s.i], s.d
}

type scroll struct {
	rs      []rune
	f       string
	i, w, d int
	b       bool
}

func (s *scroll) Frame() (string, time.Duration) {
	f := fmt.Sprintf(s.f, string(s.rs[s.i:s.i+s.w]))
	s.i += s.d
	switch {
	case s.b && (s.i < 0 || s.i > len(s.rs)-s.w):
		s.d = -s.d
		s.i += s.d
	case s.i < 0:
		s.i = len(s.rs) - s.w
	case s.i > len(s.rs)-s.w:
		s.i = 0
	}
	return f, 90 * time.Millisecond
}

// NewRtL creates a Style for a text scroller with the specified width
// and format. It scrolls from right to left.
func NewRtL(width int, format, text string) Style {
	return &scroll{
		rs: []rune(fmt.Sprintf(fmt.Sprintf("%%%vv%%v%%%[1]v[1]v", width), "", text)),
		f:  format,
		w:  width,
		d:  1,
	}
}

// NewLtR creates a Style for a text scroller with the specified
// width and format. It scrolls from left to right.
func NewLtR(width int, format, text string) Style {
	return &scroll{
		rs: []rune(fmt.Sprintf(fmt.Sprintf("%%%vv%%v%%%[1]v[1]v", width), "", text)),
		f:  format,
		w:  width,
		i:  width + len([]rune(text)),
		d:  -1,
	}
}

// NewBounce creates a Style with some text that bounces back and forth.
func NewBounce(width int, format, text string) Style {
	p := width - len([]rune(text))
	if p < 0 {
		p = 0
	}
	return &scroll{
		rs: []rune(fmt.Sprintf(fmt.Sprintf("%%%vv%%v%%%[1]v[1]v", p), "", text)),
		f:  format,
		w:  width,
		d:  1,
		b:  true,
	}
}

type wave struct{}

func (_ wave) Frame() (string, time.Duration) {
	bs := map[int]string{
		0: "⎺", 1: "⎻", 2: "⎼", 11: "⎻", 12: "⎼", 21: "⎻", 22: "⎼", 3: "⎽",
	}

	ds := map[time.Duration]struct{}{
		90 * time.Millisecond: {}, 80 * time.Millisecond: {},
		100 * time.Millisecond: {}, 70 * time.Millisecond: {},
	}

	var (
		f string
		d time.Duration
	)
	for i := 0; i < 5; i++ {
		for _, b := range bs {
			f += b
			break
		}
	}
	for d = range ds {
		break
	}
	return f, d
}
