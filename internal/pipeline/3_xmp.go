package pipeline

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func txtToXMP(logger *slog.Logger, progress *ProgressBar, txtDir string, xmpDir string, txtFiles []string, geoByBase map[string]GeoMetadata, dryRun bool) error {
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

		keywords, caption := parseKeywordsCaptionTxt(string(content))
		geo := geoByBase[base]
		xmp := buildXMP(keywords, caption, geo)

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

func parseKeywordsCaptionTxt(content string) ([]string, string) {
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

func buildXMP(keywords []string, caption string, geo GeoMetadata) string {
	captionEscaped := xmlEscape(caption)
	var keywordItems strings.Builder
	for _, kw := range keywords {
		keywordItems.WriteString("<rdf:li>")
		keywordItems.WriteString(xmlEscape(kw))
		keywordItems.WriteString("</rdf:li>")
	}

	geoFields := ""
	if geo.Valid {
		cityEscaped := xmlEscape(geo.City)
		sublocationEscaped := xmlEscape(geo.Sublocation)
		stateEscaped := xmlEscape(geo.State)
		countryEscaped := xmlEscape(geo.Country)
		countryCodeEscaped := xmlEscape(geo.CountryCode)
		lat := formatXMPGPSCoord(geo.Latitude, true)
		lon := formatXMPGPSCoord(geo.Longitude, false)
		altitudeRef := "0"
		if geo.Altitude < 0 {
			altitudeRef = "1"
		}
		geoFields = fmt.Sprintf(`
    <Iptc4xmpCore:Location>%s</Iptc4xmpCore:Location>
    <photoshop:City>%s</photoshop:City>
    <photoshop:State>%s</photoshop:State>
    <photoshop:Country>%s</photoshop:Country>
    <Iptc4xmpCore:CountryCode>%s</Iptc4xmpCore:CountryCode>
    <exif:GPSLatitude>%s</exif:GPSLatitude>
    <exif:GPSLongitude>%s</exif:GPSLongitude>
    <exif:GPSAltitude>%0.2f</exif:GPSAltitude>
    <exif:GPSAltitudeRef>%s</exif:GPSAltitudeRef>
`, sublocationEscaped, cityEscaped, stateEscaped, countryEscaped, countryCodeEscaped, lat, lon, geo.Altitude, altitudeRef)
	}

	return fmt.Sprintf(`<x:xmpmeta xmlns:x="adobe:ns:meta/">
 <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
  <rdf:Description rdf:about=""
    xmlns:dc="http://purl.org/dc/elements/1.1/"
    xmlns:lr="http://ns.adobe.com/lightroom/1.0/"
    xmlns:photoshop="http://ns.adobe.com/photoshop/1.0/"
    xmlns:exif="http://ns.adobe.com/exif/1.0/"
    xmlns:Iptc4xmpCore="http://iptc.org/std/Iptc4xmpCore/1.0/xmlns/">

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
    </dc:subject>%s

  </rdf:Description>
 </rdf:RDF>
</x:xmpmeta>
`, captionEscaped, captionEscaped, keywordItems.String(), geoFields)
}

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
