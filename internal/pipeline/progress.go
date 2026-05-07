package pipeline

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	width int
}

func NewProgressBar() *ProgressBar {
	return &ProgressBar{width: 30}
}

func (p *ProgressBar) Render(current, total int) {
	if total <= 0 {
		total = 1
	}
	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}

	filled := current * p.width / total
	empty := p.width - filled

	filledBar := strings.Repeat("#", filled)
	emptyBar := strings.Repeat("-", empty)

	green := "\033[32m"
	gray := "\033[90m"
	reset := "\033[0m"

	fmt.Printf("\r[%s%s%s%s] %d/%d", green, filledBar, gray, emptyBar, reset, current, total)
	if current == total {
		fmt.Print("\n")
	}
}
