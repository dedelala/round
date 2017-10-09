package round

import (
	"fmt"
	"time"
)

// Style is a spinner style. Any number of frames is allowed, each frame can be any length.
// The following Styles are supplied: Pipe, Block, Moon, Hearts.
type Style struct {
	Frames []string
	Rate   time.Duration
}

// Block is a good old trusty block thing. █  2580—259F Block Elements.
var Block = Style{
	[]string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█", "▇", "▆", "▅", "▄", "▃", "▂", "▁", ""},
	60 * time.Millisecond,
}

// Cylon is ominous. @ 0020—007F Basic Latin.
var Cylon = NewBouncer(7, "\x1b[1m(\x1b[31m%v\x1b[0;1m)\x1b[0m", "@")

// Hearts is clearly the best style! 💜 1F300—1F5FF Misc Symbols and Pictographs.
var Hearts = Style{
	[]string{"💖💛💚💙💜", "💜💖💛💚💙", "💙💜💖💛💚", "💚💙💜💖💛", "💛💚💙💜💖"},
	90 * time.Millisecond,
}

// Moon is a series of moon phases. 🌕 1F300—1F5FF Misc Symbols and Pictographs.
var Moon = Style{[]string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}, 90 * time.Millisecond}

// Pipe is guaranteed to work. | 0020—007F Basic Latin.
var Pipe = Style{[]string{"|", "/", "-", "\\"}, 60 * time.Millisecond}

// NewScroller creates a Style for a text scroller with the specified width
// and format. It scrolls from right to left.
func NewScroller(width int, format, text string) Style {
	r := []rune(fmt.Sprintf(fmt.Sprintf("%%%vv%%v%%%[1]v[1]v", width), "", text))
	s := Style{[]string{}, 90 * time.Millisecond}
	for i := 0; i < len(r)-width; i++ {
		s.Frames = append(s.Frames, fmt.Sprintf(format, string(r[i:i+width])))
	}
	return s
}

// NewInvertedScroller creates a Style for a text scroller with the specified
// width and format. It scrolls from left to right.
func NewInvertedScroller(width int, format, text string) Style {
	r := []rune(fmt.Sprintf(fmt.Sprintf("%%%vv%%v%%%[1]v[1]v", width), "", text))
	s := Style{[]string{}, 90 * time.Millisecond}
	for i := len(r) - width; i >= 0; i-- {
		s.Frames = append(s.Frames, fmt.Sprintf(format, string(r[i:i+width])))
	}
	return s
}

// NewBouncer creates a Style with some text that bounces back and forth.
func NewBouncer(width int, format, text string) Style {
	a := NewScroller(width, format, text)
	b := NewInvertedScroller(width, format, text)
	a.Frames = append(a.Frames[1:len(a.Frames)-1], b.Frames[1:len(b.Frames)-1]...)
	return a
}

// NewInvertedBouncer creates a Style just like NewBouncer except it goes in
// the opposite direction at the start. For the true pedant.
func NewInvertedBouncer(width int, format, text string) Style {
	a := NewScroller(width, format, text)
	b := NewInvertedScroller(width, format, text)
	a.Frames = append(b.Frames[1:len(a.Frames)-1], a.Frames[1:len(b.Frames)-1]...)
	return a
}
