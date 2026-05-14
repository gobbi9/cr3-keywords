package pipeline

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// GeoMetadata stores geotag and resolved location data for a single image.
type GeoMetadata struct {
	Valid       bool
	Latitude    float64
	Longitude   float64
	Altitude    float64
	City        string
	Sublocation string
	State       string
	Country     string
	CountryCode string
	CapturedAt  time.Time
	TrackTime   time.Time
	TimeDelta   time.Duration
}

type gpxRoot struct {
	Tracks []gpxTrack `xml:"trk"`
}

type gpxTrack struct {
	Segments []gpxSegment `xml:"trkseg"`
}

type gpxSegment struct {
	Points []gpxPoint `xml:"trkpt"`
}

type gpxPoint struct {
	Lat  float64 `xml:"lat,attr"`
	Lon  float64 `xml:"lon,attr"`
	Ele  float64 `xml:"ele"`
	Time string  `xml:"time"`
	Name string  `xml:"name"`
}

type trackPoint struct {
	Lat  float64
	Lon  float64
	Ele  float64
	Time time.Time
	Name string
}

type reverseLocation struct {
	City        string
	Sublocation string
	State       string
	Country     string
	CountryCode string
}

type reverseGeocoder struct {
	httpClient *http.Client
	mu         sync.Mutex
	cache      map[string]reverseLocation
}

func newReverseGeocoder() *reverseGeocoder {
	return &reverseGeocoder{
		httpClient: &http.Client{Timeout: 20 * time.Second},
		cache:      make(map[string]reverseLocation),
	}
}

func (r *reverseGeocoder) Lookup(ctx context.Context, lat float64, lon float64) (reverseLocation, error) {
	key := fmt.Sprintf("%.5f,%.5f", lat, lon)

	r.mu.Lock()
	if v, ok := r.cache[key]; ok {
		r.mu.Unlock()
		return v, nil
	}
	r.mu.Unlock()

	params := url.Values{}
	params.Set("format", "jsonv2")
	params.Set("lat", fmt.Sprintf("%.7f", lat))
	params.Set("lon", fmt.Sprintf("%.7f", lon))
	params.Set("zoom", "18")
	params.Set("addressdetails", "1")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://nominatim.openstreetmap.org/reverse?"+params.Encode(), nil)
	if err != nil {
		return reverseLocation{}, err
	}
	req.Header.Set("User-Agent", "cr3-keywords/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return reverseLocation{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return reverseLocation{}, fmt.Errorf("reverse geocode returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return reverseLocation{}, err
	}

	var out struct {
		DisplayName string `json:"display_name"`
		Address     struct {
			City          string `json:"city"`
			Town          string `json:"town"`
			Village       string `json:"village"`
			Hamlet        string `json:"hamlet"`
			Municipality  string `json:"municipality"`
			State         string `json:"state"`
			Country       string `json:"country"`
			CountryCode   string `json:"country_code"`
			Road          string `json:"road"`
			Pedestrian    string `json:"pedestrian"`
			Footway       string `json:"footway"`
			Path          string `json:"path"`
			Residential   string `json:"residential"`
			Neighbourhood string `json:"neighbourhood"`
			Suburb        string `json:"suburb"`
			Attraction    string `json:"attraction"`
			Building      string `json:"building"`
			Amenity       string `json:"amenity"`
		} `json:"address"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return reverseLocation{}, err
	}

	city := firstNonEmpty(
		out.Address.City,
		out.Address.Town,
		out.Address.Village,
		out.Address.Municipality,
		out.Address.Hamlet,
	)
	sublocation := firstNonEmpty(
		out.Address.Attraction,
		out.Address.Amenity,
		out.Address.Building,
		out.Address.Road,
		out.Address.Pedestrian,
		out.Address.Footway,
		out.Address.Path,
		out.Address.Residential,
		out.Address.Neighbourhood,
		out.Address.Suburb,
	)

	if sublocation == "" {
		sublocation = out.DisplayName
	}
	v := reverseLocation{
		City:        strings.TrimSpace(city),
		Sublocation: strings.TrimSpace(sublocation),
		State:       strings.TrimSpace(out.Address.State),
		Country:     strings.TrimSpace(out.Address.Country),
		CountryCode: strings.ToUpper(strings.TrimSpace(out.Address.CountryCode)),
	}

	r.mu.Lock()
	r.cache[key] = v
	r.mu.Unlock()

	return v, nil
}

// BuildGeoMetadata matches each CR3 file against the GPX track and enriches data
// with reverse-geocoded city/sublocation.
func BuildGeoMetadata(ctx context.Context, logger *slog.Logger, cr3Files []string, gpxPath string) (map[string]GeoMetadata, error) {
	points, err := parseGPXFile(gpxPath)
	if err != nil {
		return nil, fmt.Errorf("parse GPX file: %w", err)
	}
	if len(points) == 0 {
		return nil, fmt.Errorf("GPX file contains no track points")
	}

	geocoder := newReverseGeocoder()
	geoByBase := make(map[string]GeoMetadata, len(cr3Files))

	for _, cr3 := range cr3Files {
		base := strings.TrimSuffix(filepath.Base(cr3), filepath.Ext(cr3))
		capturedAt, capErr := readCaptureTime(ctx, cr3)
		if capErr != nil {
			logger.Warn("capture timestamp unavailable, using file modification time", "file", cr3, "error", capErr)
			st, statErr := os.Stat(cr3)
			if statErr != nil {
				logger.Warn("failed to stat file for timestamp fallback", "file", cr3, "error", statErr)
				continue
			}
			capturedAt = st.ModTime()
		}

		p, ok := nearestTrackPoint(points, capturedAt)
		if !ok {
			logger.Warn("no GPX point matched image", "file", cr3, "captured_at", capturedAt.Format(time.RFC3339))
			continue
		}

		loc, locErr := geocoder.Lookup(ctx, p.Lat, p.Lon)
		if locErr != nil {
			logger.Warn("reverse geocoding failed", "file", cr3, "error", locErr)
		}

		sublocation := strings.TrimSpace(loc.Sublocation)
		if sublocation == "" {
			sublocation = strings.TrimSpace(p.Name)
		}
		city := strings.TrimSpace(loc.City)
		state := strings.TrimSpace(loc.State)
		country := strings.TrimSpace(loc.Country)
		countryCode := strings.TrimSpace(loc.CountryCode)

		timeDelta := absDuration(p.Time.Sub(capturedAt))
		geoByBase[base] = GeoMetadata{
			Valid:       true,
			Latitude:    p.Lat,
			Longitude:   p.Lon,
			Altitude:    p.Ele,
			City:        city,
			Sublocation: sublocation,
			State:       state,
			Country:     country,
			CountryCode: countryCode,
			CapturedAt:  capturedAt,
			TrackTime:   p.Time,
			TimeDelta:   timeDelta,
		}

		logger.Debug("geotag matched", "file", cr3, "captured_at", capturedAt.Format(time.RFC3339), "track_time", p.Time.Format(time.RFC3339), "time_delta", timeDelta.String(), "lat", fmt.Sprintf("%.6f", p.Lat), "lon", fmt.Sprintf("%.6f", p.Lon), "city", city, "sublocation", sublocation, "state", state, "country", country, "country_code", countryCode)
	}

	return geoByBase, nil
}

func parseGPXFile(path string) ([]trackPoint, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var root gpxRoot
	if err := xml.Unmarshal(b, &root); err != nil {
		return nil, err
	}

	points := make([]trackPoint, 0)
	for _, trk := range root.Tracks {
		for _, seg := range trk.Segments {
			for _, pt := range seg.Points {
				if strings.TrimSpace(pt.Time) == "" {
					continue
				}
				t, err := parseGPXTime(pt.Time)
				if err != nil {
					continue
				}
				points = append(points, trackPoint{
					Lat:  pt.Lat,
					Lon:  pt.Lon,
					Ele:  pt.Ele,
					Time: t,
					Name: strings.TrimSpace(pt.Name),
				})
			}
		}
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].Time.Before(points[j].Time)
	})

	return points, nil
}

func parseGPXTime(v string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, strings.TrimSpace(v)); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported GPX time: %q", v)
}

func nearestTrackPoint(points []trackPoint, target time.Time) (trackPoint, bool) {
	if len(points) == 0 {
		return trackPoint{}, false
	}
	best := points[0]
	bestDiff := absDuration(best.Time.Sub(target))

	for i := 1; i < len(points); i++ {
		d := absDuration(points[i].Time.Sub(target))
		if d < bestDiff {
			best = points[i]
			bestDiff = d
		}
	}
	return best, true
}

func readCaptureTime(ctx context.Context, cr3Path string) (time.Time, error) {
	if t, err := readCaptureTimeFromExif(cr3Path); err == nil {
		return t, nil
	}
	return readCaptureTimeWithExifTool(ctx, cr3Path)
}

func readCaptureTimeFromExif(cr3Path string) (time.Time, error) {
	f, err := os.Open(cr3Path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	t, err := x.DateTime()
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func readCaptureTimeWithExifTool(ctx context.Context, cr3Path string) (time.Time, error) {
	if _, err := exec.LookPath("exiftool"); err != nil {
		return time.Time{}, fmt.Errorf("exiftool not found in PATH")
	}

	cmd := exec.CommandContext(ctx, "exiftool", "-j", "-DateTimeOriginal", "-SubSecDateTimeOriginal", "-CreateDate", cr3Path)
	out, err := cmd.Output()
	if err != nil {
		return time.Time{}, err
	}

	var rows []map[string]any
	if err := json.Unmarshal(out, &rows); err != nil {
		return time.Time{}, err
	}
	if len(rows) == 0 {
		return time.Time{}, fmt.Errorf("empty exiftool output")
	}

	for _, key := range []string{"SubSecDateTimeOriginal", "DateTimeOriginal", "CreateDate"} {
		if raw, ok := rows[0][key]; ok {
			if s, ok := raw.(string); ok {
				if t, err := parseExifDateTime(s); err == nil {
					return t, nil
				}
			}
		}
	}

	return time.Time{}, fmt.Errorf("no usable capture timestamp in EXIF")
}

func parseExifDateTime(v string) (time.Time, error) {
	v = strings.TrimSpace(v)
	layouts := []string{
		"2006:01:02 15:04:05.999999-07:00",
		"2006:01:02 15:04:05-07:00",
		"2006:01:02 15:04:05.999999",
		"2006:01:02 15:04:05",
	}
	for _, layout := range layouts {
		if strings.Contains(layout, "-07:00") {
			if t, err := time.Parse(layout, v); err == nil {
				return t, nil
			}
			continue
		}
		if t, err := time.ParseInLocation(layout, v, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported EXIF date format: %q", v)
}

func prependGeoContext(prompt string, geo GeoMetadata) string {
	if !geo.Valid {
		return prompt
	}

	city := strings.TrimSpace(geo.City)
	sublocation := strings.TrimSpace(geo.Sublocation)
	if city == "" {
		city = "unknown city"
	}
	if sublocation == "" {
		sublocation = "unknown sublocation"
	}

	prefix := fmt.Sprintf("Location context:\n- city: %s\n- sublocation: %s\n\n", city, sublocation)
	return prefix + prompt
}

func formatXMPGPSCoord(decimal float64, isLat bool) string {
	absVal := math.Abs(decimal)
	deg := math.Floor(absVal)
	min := (absVal - deg) * 60.0
	ref := "N"
	if isLat {
		if decimal < 0 {
			ref = "S"
		}
	} else {
		ref = "E"
		if decimal < 0 {
			ref = "W"
		}
	}
	return fmt.Sprintf("%.0f,%.6f%s", deg, min, ref)
}

func absDuration(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
