// acmeaddr evaluates the address taken from the arguments and prints the
// character offsets of that address.
//
// acmeaddr resets the state of /mnt/acme/$winid/addr to the range of
// the current selection.
//
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"9fans.net/go/acme"
)

func usage() error {
	return fmt.Errorf("acmeaddr [expr]")
}

func runes(w *acme.Win, expr string) (int, int, error) {

	// The addr file must be opened for addr=dot to take effect on
	// reading.
	//
	if _, _, err := w.ReadAddr(); err != nil {
		return 0, 0, fmt.Errorf("runes failed to read addr: %w", err)
	}

	if err := w.Ctl("addr=dot\n"); err != nil {
		return 0, 0, fmt.Errorf("runes failed to addr=dot: %w", err)
	}

	if err := w.Addr("%s", expr); err != nil {
		return 0, 0, fmt.Errorf("runes failed to write to addr: %w", err)
	}

	i, j, err := w.ReadAddr()
	if err != nil {
		return 0, 0, fmt.Errorf("runes failed to read addr: %w", err)
	}

	return i, j, nil
}

func run() error {
	var expr string

	switch len(os.Args) {
	case 1:
		expr = "."
	case 2:
		expr = os.Args[1]
	default:
		return usage()
	}

	id, err := strconv.Atoi(os.Getenv("winid"))
	if err != nil {
		return fmt.Errorf("run has a non-number winid: %w", err)
	}
	if id == 0 {
		return fmt.Errorf("run has an invalid winid: %d", id)
	}

	w, err := acme.Open(id, nil)
	if err != nil {
		return fmt.Errorf("run failed to open a window: %w", err)
	}
	defer w.CloseFiles()

	i, j, err := runes(w, expr)
	if err != nil {
		return fmt.Errorf("run failed to evaluate an expression: %w", err)
	}

	if i == j {
		fmt.Printf("%d\n", i)
	} else {
		fmt.Printf("%d %d\n", i, j)
	}

	return nil
}

func main() {
	log.SetPrefix("acmeaddr: ")
	log.SetFlags(0)

	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
