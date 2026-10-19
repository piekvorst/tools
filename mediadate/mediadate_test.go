package main_test

import (
	"testing"
	"time"

	"github.com/barasher/go-exiftool"

	mediadate "github.com/piekvorst/tools/mediadate"
)

func TestTimefromexiftool(t *testing.T) {
	f := func(t *testing.T, name string, info exiftool.FileMetadata, want time.Time, wanterr bool) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			got, err := mediadate.Timefromexiftool(info)

			if wanterr {
				if err == nil {
					t.Errorf("error is expected, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if got.Format(time.RFC3339Nano) != want.Format(time.RFC3339Nano) {
					t.Errorf("unexpected: %v", got)
				}
			}
		})
	}

	t.Run("from image", func(t *testing.T) {
		f(
			t, "DateTimeOriginal+SubSecTimeOriginal+OffsetTimeOriginal is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:DateTimeOriginal":   "2024:05:20 09:42:04",
					"EXIF:SubSecTimeOriginal": 465,
					"EXIF:OffsetTimeOriginal": "-07:00",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": 463,
					"EXIF:OffsetTimeDigitized": "-05:00",

					"EXIF:OffsetTime": "-04:00",
				},
			},
			time.Date(2024, 5, 20, 9, 42, 4, 465_000_000, time.FixedZone("PDT", -7*60*60)), false,
		)

		f(
			t, "DateTimeOriginal+OffsetTimeOriginal is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:DateTimeOriginal":   "2024:05:20 09:42:04",
					"EXIF:SubSecTimeOriginal": -1,
					"EXIF:OffsetTimeOriginal": "-07:00",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": 463,
					"EXIF:OffsetTimeDigitized": "-05:00",

					"EXIF:OffsetTime": "-04:00",
				},
			},
			time.Date(2024, 5, 20, 9, 42, 4, 0, time.FixedZone("PDT", -7*60*60)), false,
		)

		f(
			t, "DateTimeOriginal+OffsetTime is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:DateTimeOriginal":   "2024:05:20 09:42:04",
					"EXIF:SubSecTimeOriginal": -1,
					"EXIF:OffsetTimeOriginal": "bogus",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": 463,
					"EXIF:OffsetTimeDigitized": "-05:00",

					"EXIF:OffsetTime": "-04:00",
				},
			},
			time.Date(2024, 5, 20, 9, 42, 4, 0, time.FixedZone("AST", -4*60*60)), false,
		)

		f(
			t, "DateTimeOriginal+UTC is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:DateTimeOriginal":   "2024:05:20 09:42:04",
					"EXIF:SubSecTimeOriginal": -1,
					"EXIF:OffsetTimeOriginal": "bogus",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": 463,
					"EXIF:OffsetTimeDigitized": "-05:00",

					"EXIF:OffsetTime": "bogus",
				},
			},
			time.Date(2024, 5, 20, 9, 42, 4, 0, time.UTC), false,
		)

		f(
			t, "CreateDate+SubSecTimeDigitized+OffsetTimeDigitized is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:DateTimeOriginal":   "0000:00:00 00:00:00",
					"EXIF:SubSecTimeOriginal": -1,
					"EXIF:OffsetTimeOriginal": "-07:00",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": 463,
					"EXIF:OffsetTimeDigitized": "-05:00",

					"EXIF:OffsetTime": "-04:00",
				},
			},
			time.Date(2022, 3, 18, 7, 40, 2, 463_000_000, time.FixedZone("CDT", -5*60*60)), false,
		)

		f(
			t, "CreateDate+OffsetTimeDigitized is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": -1,
					"EXIF:OffsetTimeDigitized": "-05:00",

					"EXIF:OffsetTime": "-04:00",
				},
			},
			time.Date(2022, 3, 18, 7, 40, 2, 0, time.FixedZone("CDT", -5*60*60)), false,
		)

		f(
			t, "CreateDate+OffsetTime is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": -1,
					"EXIF:OffsetTimeDigitized": "bogus",

					"EXIF:OffsetTime": "-04:00",
				},
			},
			time.Date(2022, 3, 18, 7, 40, 2, 0, time.FixedZone("AST", -4*60*60)), false,
		)

		f(
			t, "CreateDate+UTC is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:CreateDate":          "2022:03:18 07:40:02",
					"EXIF:SubSecTimeDigitized": -1,
					"EXIF:OffsetTimeDigitized": "bogus",

					"EXIF:OffsetTime": "bogus",
				},
			},
			time.Date(2022, 3, 18, 7, 40, 2, 0, time.UTC), false,
		)

		f(
			t, "error",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "image/heic",

					"EXIF:DateTimeOriginal":   "0000:00:00 00:00:00",
					"EXIF:SubSecTimeOriginal": 465,
					"EXIF:OffsetTimeOriginal": "-07:00",

					"EXIF:CreateDate":          "0000:00:00 00:00:00",
					"EXIF:SubSecTimeDigitized": 463,
					"EXIF:OffsetTimeDigitized": "-05:00",

					"EXIF:OffsetTime": "-04:00",
				},
			},
			time.Time{}, true,
		)
	})

	t.Run("from video", func(t *testing.T) {
		f(
			t, "CreateDate is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "video/quicktime",

					"QuickTime:CreateDate":      "2024:05:20 09:42:04",
					"QuickTime:MediaCreateDate": "2023:04:19 08:41:03",
					"QuickTime:TrackCreateDate": "2022:03:18 07:40:02",
				},
			},
			time.Date(2024, 5, 20, 9, 42, 4, 0, time.UTC), false,
		)

		f(
			t, "MediaCreateDate is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "video/quicktime",

					"QuickTime:CreateDate":      "0000:00:00 00:00:00",
					"QuickTime:MediaCreateDate": "2023:04:19 08:41:03",
					"QuickTime:TrackCreateDate": "2022:03:18 07:40:02",
				},
			},
			time.Date(2023, 4, 19, 8, 41, 3, 0, time.UTC), false,
		)

		f(
			t, "TrackCreateDate is used",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "video/quicktime",

					"QuickTime:CreateDate":      "0000:00:00 00:00:00",
					"QuickTime:MediaCreateDate": "0000:00:00 00:00:00",
					"QuickTime:TrackCreateDate": "2022:03:18 07:40:02",
				},
			},
			time.Date(2022, 3, 18, 7, 40, 2, 0, time.UTC), false,
		)

		f(
			t, "error",
			exiftool.FileMetadata{
				Fields: map[string]any{
					"File:MIMEType": "video/quicktime",

					"QuickTime:CreateDate":      "0000:00:00 00:00:00",
					"QuickTime:MediaCreateDate": "0000:00:00 00:00:00",
					"QuickTime:TrackCreateDate": "0000:00:00 00:00:00",
				},
			},
			time.Time{}, true,
		)
	})
}

// TestNewDateTime assumes that NewDateTime relies on time.Parse
// internally. Thus, it does not test what time.Parse guarantees
// already.
func TestNewDateTime(t *testing.T) {
	f := func(t *testing.T, name string, v any, want time.Time, wanterr bool) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			got, err := mediadate.NewDateTime(v)

			if wanterr {
				if err == nil {
					t.Errorf("error is expected, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if !want.Equal(got) {
					t.Errorf("unexpected: %v", got)
				}
			}
		})
	}

	f(t, "happy path", "2024:05:20 09:42:04", time.Date(2024, 5, 20, 9, 42, 4, 0, time.UTC), false)
	f(t, "nil", nil, time.Time{}, true)
	f(t, "non-string", 1716198124, time.Time{}, true)
	f(t, "layout contract: empty string", "", time.Time{}, true)
	f(t, "layout contract: date only, no time", "2024:05:20", time.Time{}, true)
	f(t, "layout contract: trailing garbage", "2024:05:20 09:42:04Z", time.Time{}, true)
}

func TestNewLocation(t *testing.T) {
	f := func(t *testing.T, name string, v any, want *time.Location, wanterr bool) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			got, err := mediadate.NewLocation(v)

			if wanterr {
				if err == nil {
					t.Errorf("error is expected, got none")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				if want.String() != got.String() {
					t.Errorf("unexpected: %v", got)
				}
			}
		})
	}

	t.Run("zero offset", func(t *testing.T) {

		// why not time.UTC everywhere: do not erase the distinction

		f(t, "Z", "Z", time.UTC, false)
		f(t, "positive zero", "+00:00", time.FixedZone("+00:00", 0), false)
		f(t, "negative zero", "-00:00", time.FixedZone("-00:00", 0), false)
	})

	t.Run("±hh:mm", func(t *testing.T) {
		f(t, "negative", "-05:00", time.FixedZone("-05:00", -5*60*60), false)
		f(t, "positive", "+05:00", time.FixedZone("+05:00", 5*60*60), false)

		f(t, "half-hour", "-05:30", time.FixedZone("-05:30", -(5*60*60+30*60)), false)
		f(t, "quarter-hour", "-05:45", time.FixedZone("-05:45", -(5*60*60+45*60)), false)

		f(t, "real min", "-12:00", time.FixedZone("-12:00", -12*60*60), false)
		f(t, "real max", "+14:00", time.FixedZone("+14:00", 14*60*60), false)
	})

	t.Run("±hhmm", func(t *testing.T) {
		f(t, "negative", "-0500", time.FixedZone("-0500", -5*60*60), false)
		f(t, "positive", "+0500", time.FixedZone("+0500", 5*60*60), false)
	})

	t.Run("±hh", func(t *testing.T) {
		f(t, "negative", "-05", time.FixedZone("-05", -5*60*60), false)
		f(t, "positive", "+05", time.FixedZone("+05", 5*60*60), false)
	})

	t.Run("boundaries", func(t *testing.T) {
		f(t, "min", "-23:59", time.FixedZone("-23:59", -(23*60*60+59*60)), false)
		f(t, "max", "+23:59", time.FixedZone("+23:59", 23*60*60+59*60), false)

		f(t, "hour out of range", "-24:00", nil, true)
		f(t, "minute out of range", "-00:60", nil, true)
		f(t, "hour and minute out of range", "-24:60", nil, true)
	})

	t.Run("malformed", func(t *testing.T) {
		f(t, "empty", "", nil, true)
		f(t, "sign only", "-", nil, true)
		f(t, "lowercase z", "z", nil, true)
		f(t, "UTC name", "UTC", nil, true)
		f(t, "GMT name", "GMT", nil, true)

		f(t, "space after sign", "- 05:00", nil, true)
		f(t, "leading space", " -05:00", nil, true)
		f(t, "trailing space", "-05:00 ", nil, true)

		f(t, "single-digit hour ±hh:mm", "-5:00", nil, true)
		f(t, "single-digit hour ±hhmm", "-500", nil, true)
		f(t, "single-digit hour ±hh", "-5", nil, true)

		f(t, "single-digit minute ±hh:mm", "-05:0", nil, true)
		f(t, "single-digit minute ±hhmm", "-050", nil, true)

		f(t, "extra digit ±hh:mm", "-05:000", nil, true)
		f(t, "extra digit ±hhmm", "-05000", nil, true)
		f(t, "extra digit ±hh", "-050", nil, true)

		f(t, "dot separator", "-05.00", nil, true)
		f(t, "trailing colon", "-05:", nil, true)
		f(t, "offset with seconds", "-05:00:00", nil, true)

		f(t, "non-digit hours ±hh:mm", "-OO:OO", nil, true)
		f(t, "non-digit hours ±hhmm", "-OOOO", nil, true)
		f(t, "non-digit hours ±hh", "-OO", nil, true)

		f(t, "non-digit minutes ±hh:mm", "-05:OO", nil, true)
		f(t, "non-digit minutes ±hhmm", "-05OO", nil, true)

		f(t, "trailing garbage", "-05:00x", nil, true)
	})
}
