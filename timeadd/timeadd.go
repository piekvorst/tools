package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type parsedArgs struct {
	Seconds, Minutes, Hours time.Duration
	Days, Months, Years     int

	Time time.Time
}

type unixDateFlag time.Time

func (f *unixDateFlag) String() string {
	if f == nil {
		return ""
	}

	return time.Time(*f).Format(time.UnixDate)
}

func (f *unixDateFlag) Set(s string) error {
	parsed, err := time.Parse(time.UnixDate, s)
	if err != nil {
		return fmt.Errorf("unixDateFlag.Set: must match time.UnixDate: %w", err)
	}

	*f = unixDateFlag(parsed)
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "timeadd -m -6 Mon Jan  2 15:04:05 MST 2006")
	os.Exit(2)
}

func parseArgs(args []string) parsedArgs {
	var parsed parsedArgs

	fs := flag.NewFlagSet("timeadd", flag.ExitOnError)
	fs.Usage = usage

	var seconds, minutes, hours int
	fs.IntVar(&seconds, "s", 0, "seconds")
	fs.IntVar(&minutes, "m", 0, "minutes")
	fs.IntVar(&hours, "h", 0, "hours")

	fs.IntVar(&parsed.Days, "D", 0, "days")
	fs.IntVar(&parsed.Months, "M", 0, "months")
	fs.IntVar(&parsed.Years, "Y", 0, "years")

	fs.Parse(args)

	parsed.Seconds = time.Duration(seconds) * time.Second
	parsed.Minutes = time.Duration(minutes) * time.Minute
	parsed.Hours = time.Duration(hours) * time.Hour

	if len(fs.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "timeadd: date is absent")
		usage()
	}

	if err := (*unixDateFlag)(&parsed.Time).Set(strings.Join(fs.Args(), " ")); err != nil {
		fmt.Fprintf(os.Stderr, "timeadd: %s\n", err)
		usage()
	}

	return parsed
}

func run() error {
	args := parseArgs(os.Args[1:])

	fmt.Println(
		args.Time.
			AddDate(args.Years, args.Months, args.Days).
			Add(args.Seconds + args.Minutes + args.Hours).
			Format(time.UnixDate),
	)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "timeadd: %s\n", err)
		os.Exit(1)
	}
}
