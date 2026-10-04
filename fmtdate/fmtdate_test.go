package main_test

import (
	"flag"
	"io"
	"strings"
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

func TestEnumFlag(t *testing.T) {
	f := func(
		t *testing.T,
		name string,
		defaultvalue string,
		allowed []string,
		arguments []string,
		want string,
		wanterr string,
	) {
		t.Helper()

		t.Run(name, func(t *testing.T) {
			t.Helper()

			var got string

			fs := flag.NewFlagSet("", flag.ContinueOnError)
			fs.SetOutput(io.Discard)

			fmtdate.EnumFlagFS(fs, &got, "f", defaultvalue, allowed)

			err := fs.Parse(arguments)
			if wanterr == "" && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if wanterr != "" && (err == nil || !strings.Contains(err.Error(), wanterr)) {
				t.Fatalf("expected error containing wanterr, got: %v", err)
			}
			if got != want {
				t.Fatalf("unexpected: %q", got)
			}
		})
	}

	f(t, "no-flag-keeps-default", "foo", []string{"foo", "bar"}, nil, "foo", "")
	f(t, "default-not-validated", "baz", []string{"foo", "bar"}, nil, "baz", "")
	f(t, "lowercase-argument", "foo", []string{"foo", "bar"}, []string{"-f", "Bar"}, "bar", "")
	f(t, "trim-argument", "foo", []string{"foo", "bar"}, []string{"-f", " bar "}, "bar", "")
	f(t, "lowercase-allowed", "foo", []string{"Foo", "Bar"}, []string{"-f", "bar"}, "bar", "")
	f(t, "trim-allowed", "foo", []string{" foo ", " bar "}, []string{"-f", "bar"}, "bar", "")
	f(t, "rejects-invalid-value", "foo", []string{"foo", "bar"}, []string{"-f", "baz"}, "foo", "must be one of")
	f(t, "empty-allowed", "foo", nil, []string{"-f", "foo"}, "foo", "must be one of")
}
