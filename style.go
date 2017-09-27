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

// Pipe is guaranteed to work. | 0020—007F Basic Latin.
var Pipe = Style{[]string{"|", "/", "-", "\\"}, 60 * time.Millisecond}

// Moon is a series of moon phases. 🌕 1F300—1F5FF Misc Symbols and Pictographs.
var Moon = Style{[]string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}, 90 * time.Millisecond}

// Block is a good old trusty block thing. █  2580—259F Block Elements.
var Block = Style{
	[]string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█", "▇", "▆", "▅", "▄", "▃", "▂", "▁", ""},
	60 * time.Millisecond,
}

// Hearts is clearly the best style! 💜 1F300—1F5FF Misc Symbols and Pictographs.
var Hearts = Style{
	[]string{"💖💛💚💙💜", "💜💖💛💚💙", "💙💜💖💛💚", "💚💙💜💖💛", "💛💚💙💜💖"},
	90 * time.Millisecond,
}

// NewScroller creates a Style for a text scroller with the specified width and format.
func NewScroller(width int, format, text string) Style {
	text = fmt.Sprintf(fmt.Sprintf("%%%vv%%v%%%[1]v[1]v", width), "", text)
	s := Style{[]string{}, 90 * time.Millisecond}
	for i := 0; i < len(text)-width; i++ {
		s.Frames = append(s.Frames, fmt.Sprintf(format, text[i:i+width]))
	}
	return s
}
