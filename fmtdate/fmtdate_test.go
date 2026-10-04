package main_test

import (
	"testing"
	"time"

	fmtdate "github.com/piekvorst/tools/fmtdate"
)

var (
	CDT  = time.FixedZone("CDT", -5*60*60)      // UTC-5
	ACST = time.FixedZone("ACST", 9*3600+30*60) // UTC+9:30
)

func TestISO9075(t *testing.T) {
	f := func(t *testing.T, name string, tm time.Time, want string) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			if got := tm.Format(fmtdate.ISO9075); got != want {
				t.Errorf("unexpected: %v", got)
			}
		})
	}

	f(
		t,
		"negative-tz",
		time.Date(2027, 2, 13, 14, 25, 26, 0, CDT),
		"2027-02-13 14:25:26-05:00",
	)
	f(
		t,
		"utc",
		time.Date(2027, 2, 13, 14, 25, 26, 0, time.UTC),
		"2027-02-13 14:25:26+00:00",
	)
	f(
		t,
		"half-hour-offset",
		time.Date(2027, 2, 13, 14, 25, 26, 0, ACST),
		"2027-02-13 14:25:26+09:30",
	)
}

func TestFormatWords(t *testing.T) {
	f := func(t *testing.T, name string, tm time.Time, want string) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			if got := fmtdate.FormatWords(tm); got != want {
				t.Errorf("unexpected: %v", got)
			}
		})
	}

	f(
		t,
		"negative-half-hour-tz",
		time.Date(2027, 2, 13, 14, 25, 26, 0, time.FixedZone("UTC-00:30", -(30*60))),
		"2027-02-13 14:25:26 - 00:30",
	)
	f(
		t,
		"utc",
		time.Date(2027, 2, 13, 14, 25, 26, 0, time.UTC),
		"2027-02-13 14:25:26 + 00:00",
	)
	f(
		t,
		"half-hour-offset",
		time.Date(2027, 2, 13, 14, 25, 26, 0, ACST),
		"2027-02-13 14:25:26 + 09:30",
	)
}
