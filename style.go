package round

import (
	"fmt"
	"time"
)

// Predefined styles!
var (
	Pipe  = Style{[]string{"|", "/", "-", "\\"}, 60 * time.Millisecond}
	Moon  = Style{[]string{"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"}, 90 * time.Millisecond}
	Block = Style{
		[]string{"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█", "▇", "▆", "▅", "▄", "▃", "▂", "▁", ""},
		60 * time.Millisecond,
	}
	Hearts = Style{
		[]string{"❤️ 💛 💚 💙 💜 ", "💜 ❤️ 💛 💚 💙 ", "💙 💜 ❤️ 💛 💚 ", "💚 💙 💜 ❤️ 💛 ", "💛 💚 💙 💜 ❤️ "},
		90 * time.Millisecond,
	}
)

// Style is a spinner style. Any number of frames is allowed,
// and each frame need not be the same length.
type Style struct {
	Frames []string
	Rate   time.Duration
}

// NewScroller creates a Style for a text scroller.
func NewScroller(width int, format, text string) Style {
	text = fmt.Sprintf(fmt.Sprintf("%%%vv%%v%%%[1]v[1]v", width), "", text)
	s := Style{[]string{}, 90 * time.Millisecond}
	for i := 0; i < len(text)-width; i++ {
		s.Frames = append(s.Frames, fmt.Sprintf(format, text[i:i+width]))
	}
	return s
}
