package pipeline

import (
	"path/filepath"
	"strings"
)

func xmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}

func folderNameFromPath(p string) string {
	clean := filepath.Clean(p)
	return filepath.Base(clean)
}
