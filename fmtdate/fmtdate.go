/*
fmtdate prints the date, time, and timezone.

The flags are:

	-f format
	    Possible values are: words (default), iso9075, rfc3339.

	-u
	    Report Greenwich Mean Time (GMT) rather than local time.

# Notes for shell usage

The default format "words" is designed specifically to be friendly
for simple text manipulation. In particular, parsing is as easy
as:

    fs = '-:
     	' # newline, space, and tab
    now = `$fs{fmtdate}

Both ISO 9075 and RFC 3339 attach the timezone directly to the
time, which makes it impossible to parse a string in a simple
word-splitting way (unless the time zone is UTC). In this case, a
long regular expression is required, which makes the code untidy
and excessively verbose.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"
)

type FormatFlag string

const (
	WordsFlag FormatFlag = "words"
	ISO9075Flag FormatFlag = "iso9075"
	RFC3339Flag FormatFlag = "rfc3339"
)

const ISO9075 = "2006-01-02 15:04:05-07:00"

func Join[T ~string](elems []T, sep string) string {
	if len(elems) == 0 {
		return ""
	}

	buf := make([]string, 0, len(elems))
	for _, v := range elems {
		buf = append(buf, string(v))
	}

	return strings.Join(buf, sep)
}

func EnumFlag[T ~string](p *T, name string, value T, allowed []T) {
	*p = value

	flag.Func(name, "", func(raw string) error {
		clean := T(strings.ToLower(strings.TrimSpace(raw)))
		if !slices.Contains(allowed, clean) {
			return fmt.Errorf("must be one of %s", Join(allowed, " "))
		}

		*p = clean

		return nil
	})
}

func FormatWords(t time.Time) string {
	type tzwords struct {
		Sign rune
		Hours int
		Minutes int
	}

	newtzwords := func(seconds int) (tz tzwords) {
		tz.Hours = seconds / 3600
		tz.Minutes = (seconds % 3600) / 60

		// +00:00 is more compatible for representing UTC.
		//
		if seconds >= 0 {
			tz.Sign = '+'
		} else {
			tz.Sign = '-'
			tz.Hours = -tz.Hours
			tz.Minutes = -tz.Minutes
		}

		return tz
	}

	_, seconds := t.Zone()
	tz := newtzwords(seconds)

	return fmt.Sprintf(
		"%s %c %02d:%02d",
		t.Format("2006-01-02 15:04:05"),
		tz.Sign,
		tz.Hours,
		tz.Minutes,
	)
}

func main() {
	log.SetPrefix("fmtdate: ")
	log.SetFlags(0)

	var format FormatFlag
	var gmt bool

	EnumFlag(&format, "f", WordsFlag, []FormatFlag{WordsFlag, ISO9075Flag, RFC3339Flag})
	flag.BoolVar(&gmt, "u", false, "")

	flag.Parse()

	var now time.Time
	if !gmt {
		now = time.Now()
	} else {
		now = time.Now().UTC()
	}

	var output string
	switch format {
	case WordsFlag:
		output = FormatWords(now)
	case ISO9075Flag:
		output = now.Format(ISO9075)
	case RFC3339Flag:
		output = now.Format(time.RFC3339)
	}

	// A safeguard against adding a new format and forgetting to
	// handle it.
	//
	if output == "" {
		log.Fatalf("unknown format: %s", format)
	}

	fmt.Println(output)
}
