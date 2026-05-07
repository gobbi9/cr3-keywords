package pipeline

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func txtToXMP(logger *slog.Logger, progress *ProgressBar, txtDir string, xmpDir string, txtFiles []string, dryRun bool) error {
	if len(txtFiles) == 0 {
		return fmt.Errorf("no TXT files to convert")
	}

	total := len(txtFiles)
	for idx, txt := range txtFiles {
		base := strings.TrimSuffix(filepath.Base(txt), filepath.Ext(txt))
		xmpFile := filepath.Join(xmpDir, base+".xmp")

		if _, err := os.Stat(xmpFile); err == nil {
			logger.Debug("skipping TXT because XMP exists", "file", txt)
			progress.Render(idx+1, total)
			continue
		}

		if dryRun {
			logger.Debug("[dry-run] would write XMP", "from", txt, "to", xmpFile)
			progress.Render(idx+1, total)
			continue
		}

		content, err := os.ReadFile(txt)
		if err != nil {
			logger.Error("failed to read TXT", "file", txt, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		keywords, caption := parseCaptionTXT(string(content))
		xmp := buildXMP(keywords, caption)

		if err := os.WriteFile(xmpFile, []byte(xmp), 0o644); err != nil {
			logger.Error("failed to write XMP", "file", xmpFile, "error", err)
			progress.Render(idx+1, total)
			continue
		}

		progress.Render(idx+1, total)
	}

	_ = txtDir
	return nil
}

func parseCaptionTXT(content string) ([]string, string) {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return nil, ""
	}

	first := strings.TrimSpace(lines[0])
	parts := strings.Split(first, ",")
	keywords := make([]string, 0, len(parts))
	for _, p := range parts {
		k := strings.TrimSpace(p)
		if k != "" {
			keywords = append(keywords, k)
		}
	}

	caption := ""
	if len(lines) > 1 {
		for i := 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "" {
				if i+1 < len(lines) {
					caption = strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
				}
				break
			}
		}
		if caption == "" {
			caption = strings.TrimSpace(strings.Join(lines[1:], "\n"))
		}
	}

	return keywords, caption
}

func buildXMP(keywords []string, caption string) string {
	captionEscaped := xmlEscape(caption)
	var keywordItems strings.Builder
	for _, kw := range keywords {
		keywordItems.WriteString("<rdf:li>")
		keywordItems.WriteString(xmlEscape(kw))
		keywordItems.WriteString("</rdf:li>")
	}

	return fmt.Sprintf(`<x:xmpmeta xmlns:x="adobe:ns:meta/">
 <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description rdf:about=""
    xmlns:dc="http://purl.org/dc/elements/1.1/"
    xmlns:lr="http://ns.adobe.com/lightroom/1.0/">

    <lr:description>%s</lr:description>

    <dc:description>
      <rdf:Alt>
        <rdf:li xml:lang="x-default">%s</rdf:li>
      </rdf:Alt>
    </dc:description>

    <dc:subject>
      <rdf:Bag>
        %s
      </rdf:Bag>
    </dc:subject>

  </rdf:Description>
 </rdf:RDF>
</x:xmpmeta>
`, captionEscaped, captionEscaped, keywordItems.String())
}
