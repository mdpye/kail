package kail

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

var (
	prefixColor = color.New(color.FgHiWhite, color.Bold)
)

type Writer interface {
	Print(event Event) error
	Fprint(w io.Writer, event Event) error
}

func NewWriter(out io.Writer) Writer {
	return &writer{out}
}

type writer struct {
	out io.Writer
}

func (w *writer) Print(ev Event) error {
	return w.Fprint(w.out, ev)
}

func (w *writer) Fprint(out io.Writer, ev Event) error {
	prefix := w.prefix(ev)

	log := ev.Log()
	lines := strings.Split(string(log), "\n")

	for _, line := range lines {
		if len(line) > 0 {
			if _, err := prefixColor.Fprint(out, prefix); err != nil {
				return err
			}
			if _, err := prefixColor.Fprint(out, ": "); err != nil {
				return err
			}

			if _, err := out.Write([]byte(line)); err != nil {
				return err
			}
			if _, err := out.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *writer) prefix(ev Event) string {
	return fmt.Sprintf("%v/%v[%v]",
		ev.Source().Namespace(),
		ev.Source().Name(),
		ev.Source().Container())
}
