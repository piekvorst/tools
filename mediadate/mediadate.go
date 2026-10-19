// mediadate file ...
package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

type exiftoolResult struct {
	Info exiftool.FileMetadata
	Err  error
}

var offsetRE = regexp.MustCompile(`^([+-])([0-9][0-9])(?::?([0-9][0-9]))?$`)

type metadata struct {
	File struct {
		MIMEType string
	}

	EXIF struct {
		DateTimeOriginal   time.Time
		SubSecTimeOriginal time.Duration
		OffsetTimeOriginal *time.Location

		CreateDate          time.Time
		SubSecTimeDigitized time.Duration
		OffsetTimeDigitized *time.Location

		OffsetTime *time.Location
	}

	QuickTime struct {
		CreateDate      time.Time
		MediaCreateDate time.Time
		TrackCreateDate time.Time
	}
}

func newDateTime(v any) (time.Time, error) {
	s, ok := v.(string)
	if !ok {
		return time.Time{}, fmt.Errorf("cannot convert %T to time.Time", v)
	}

	t, err := time.Parse("2006:01:02 15:04:05", s)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

func newLocation(v any) (*time.Location, error) {
	atoi := func(s string) int { return int(s[0]-'0')*10 + int(s[1]-'0') }

	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("unsupported offset type %T", v)
	}
	if s == "Z" {
		return time.UTC, nil
	}

	m := offsetRE.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("invalid offset: %q", s)
	}

	sign := m[1]
	hh := atoi(m[2])
	var mm int
	if m[3] != "" {
		mm = atoi(m[3])
	}
	if hh > 23 {
		return nil, fmt.Errorf("hour out of range in %q", s)
	}
	if mm > 59 {
		return nil, fmt.Errorf("minute out of range in %q", s)
	}

	off := hh*3600 + mm*60
	if sign == "-" {
		off = -off
	}
	return time.FixedZone(s, off), nil
}

func newmetadata(info exiftool.FileMetadata) *metadata {
	parsedatetime := func(k string) time.Time {
		raw, err := info.GetString(k)
		if err != nil {
			return time.Time{}
		}

		t, _ := newDateTime(strings.TrimSpace(raw))
		return t
	}

	parsesubsec := func(k string) time.Duration {
		raw, err := info.GetString(k)
		if err != nil {
			return 0
		}

		d, err := time.ParseDuration(fmt.Sprintf("0.%ss", strings.TrimSpace(raw)))
		if err != nil {
			return 0
		}

		return d
	}

	parselocation := func(k string) *time.Location {
		raw, err := info.GetString(k)
		if err != nil {
			return nil
		}

		loc, _ := newLocation(strings.TrimSpace(raw))
		return loc
	}

	var parsed metadata

	raw, _ := info.GetString("File:MIMEType")
	parsed.File.MIMEType = strings.ToLower(strings.TrimSpace(raw))

	parsed.EXIF.DateTimeOriginal = parsedatetime("EXIF:DateTimeOriginal")
	parsed.EXIF.SubSecTimeOriginal = parsesubsec("EXIF:SubSecTimeOriginal")
	parsed.EXIF.OffsetTimeOriginal = parselocation("EXIF:OffsetTimeOriginal")

	parsed.EXIF.CreateDate = parsedatetime("EXIF:CreateDate")
	parsed.EXIF.SubSecTimeDigitized = parsesubsec("EXIF:SubSecTimeDigitized")
	parsed.EXIF.OffsetTimeDigitized = parselocation("EXIF:OffsetTimeDigitized")

	parsed.EXIF.OffsetTime = parselocation("EXIF:OffsetTime")

	parsed.QuickTime.CreateDate = parsedatetime("QuickTime:CreateDate")
	parsed.QuickTime.MediaCreateDate = parsedatetime("QuickTime:MediaCreateDate")
	parsed.QuickTime.TrackCreateDate = parsedatetime("QuickTime:TrackCreateDate")

	return &parsed
}

// The following priority is used:
//
// . DateTimeOriginal + (SubSecTimeOriginal ?? 0) + (OffsetTimeOriginal ?? OffsetTime ?? UTC)
//
// . CreateDate + (SubSecTimeDigitized ?? 0) + (OffsetTimeDigitized ?? OffsetTime ?? UTC)
//
// . error
func timefromimage(info *metadata) (time.Time, error) {
	mk := func(t time.Time, nano time.Duration, locs []*time.Location) time.Time {
		if nano != 0 {
			t = t.Truncate(time.Second).Add(nano)
		}

		for _, loc := range locs {
			if loc != nil {
				t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
				break
			}
		}

		return t
	}

	if u := info.EXIF.DateTimeOriginal; !u.IsZero() {
		return mk(
			u,
			info.EXIF.SubSecTimeOriginal,
			[]*time.Location{
				info.EXIF.OffsetTimeOriginal,
				info.EXIF.OffsetTime,
				time.UTC,
			},
		), nil
	} else if u := info.EXIF.CreateDate; !u.IsZero() {
		return mk(
			u,
			info.EXIF.SubSecTimeDigitized,
			[]*time.Location{
				info.EXIF.OffsetTimeDigitized,
				info.EXIF.OffsetTime,
				time.UTC,
			},
		), nil
	} else {
		return time.Time{}, fmt.Errorf("both DateTimeOriginal and CreateDate are zero")
	}
}

// The following priority is used:
//
// . (CreateDate ?? MediaCreateDate ?? TrackCreateDate) + UTC
//
// . error
func timefromvideo(info *metadata) (time.Time, error) {
	var t time.Time
	for _, u := range []time.Time{
		info.QuickTime.CreateDate,
		info.QuickTime.MediaCreateDate,
		info.QuickTime.TrackCreateDate,
	} {
		if !u.IsZero() {
			t = u
			break
		}
	}
	if t.IsZero() {
		return time.Time{}, fmt.Errorf("CreateDate, MediaCreateDate, and TrackCreateDate are unset")
	}

	return t, nil
}

func timefromexiftool(info exiftool.FileMetadata) (time.Time, error) {
	parsed := newmetadata(info)

	switch {
	case strings.HasPrefix(parsed.File.MIMEType, "image/"):
		return timefromimage(parsed)
	case strings.HasPrefix(parsed.File.MIMEType, "video/"):
		return timefromvideo(parsed)
	default:
		return time.Time{}, fmt.Errorf("unsupported mime type: %q", parsed.File.MIMEType)
	}
}

func extractMetadataStream(ctx context.Context, et *exiftool.Exiftool, paths []string) <-chan exiftoolResult {
	out := make(chan exiftoolResult)

	go func() {
		defer close(out)

		send := func(result exiftoolResult) bool {
			select {
			case out <- result:
				return true
			case <-ctx.Done():
				return false
			}
		}

		for _, p := range paths {
			raw := et.ExtractMetadata(p)

			// enforce the implicit contract
			if len(raw) != 1 {
				send(exiftoolResult{Err: fmt.Errorf("unexpected number of return values for %s: %v", p, len(raw))})
				continue
			}

			if raw[0].Err != nil {
				send(exiftoolResult{Err: fmt.Errorf("failed to extract metadata for %s: %w", raw[0].File, raw[0].Err)})
				continue
			}

			if !send(exiftoolResult{Info: raw[0]}) {
				return
			}
		}
	}()

	return out
}

func main() {
	ctx := context.Background()

	et, err := exiftool.NewExiftool(
		exiftool.NoPrintConversion(),
		exiftool.PrintGroupNames("0"),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize exiftool: %v\n", err)
		os.Exit(1)
	}
	defer et.Close()

	paths := os.Args[1:]

	for msg := range extractMetadataStream(ctx, et, paths) {
		if msg.Err != nil {
			fmt.Fprintf(os.Stderr, "failed to extract metadata: %v\n", msg.Err)
			continue
		}

		t, err := timefromexiftool(msg.Info)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to derive time for %s: %v\n", msg.Info.File, err)
			continue
		}

		fmt.Printf("%s	%v\n", msg.Info.File, t.Unix())
	}
}
