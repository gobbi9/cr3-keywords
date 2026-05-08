package pipeline

import (
	"fmt"
	"strings"
)

// ProgressBar renders a single-line terminal progress bar.
type ProgressBar struct {
	width int
}

// NewProgressBar creates a ProgressBar with the default width.
func NewProgressBar() *ProgressBar {
	return &ProgressBar{width: 30}
}

// Render draws the progress bar for current and total units.
func (p *ProgressBar) Render(current, total int) {
	p.RenderWithDecorations(current, total, "", "")
}

// RenderWithSuffix draws the progress bar with trailing text.
func (p *ProgressBar) RenderWithSuffix(current, total int, suffix string) {
	p.RenderWithDecorations(current, total, "", suffix)
}

// RenderWithPrefix draws the progress bar with leading text.
func (p *ProgressBar) RenderWithPrefix(current, total int, prefix string) {
	p.RenderWithDecorations(current, total, prefix, "")
}

// RenderWithDecorations draws the progress bar with optional prefix and suffix.
// It clamps current into [0, total] and normalizes non-positive totals.
func (p *ProgressBar) RenderWithDecorations(current, total int, prefix string, suffix string) {
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

	fmt.Printf("\r%s[%s%s%s%s%s] %d/%d%s", prefix, green, filledBar, gray, emptyBar, reset, current, total, suffix)
	if current == total {
		fmt.Print("\n")
	}
}
