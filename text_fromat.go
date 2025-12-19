package logx

import (
	"fmt"
	"runtime"
	"time"
)

type TextFormat struct {
	Color        bool
	FormatTime   func(buffer *[]byte, t time.Time)
	FormatCaller func(buffer *[]byte, frame *runtime.Frame)
	FormatArgs   func(buffer *[]byte, args []any)
}

func (t *TextFormat) Format(buffer *[]byte, entry *Entry) {
	if t.Color {
		*buffer = append(*buffer, entry.Level.Color()...)
		*buffer = append(*buffer, ' ')
		if t.FormatTime != nil {
			t.FormatTime(buffer, entry.Time)
		} else {
			FormatTime(buffer, entry.Time)
		}
	} else {
		if t.FormatTime != nil {
			t.FormatTime(buffer, entry.Time)
		} else {
			FormatTime(buffer, entry.Time)
		}
		*buffer = append(*buffer, ' ')
		*buffer = append(*buffer, entry.Level.String()...)
	}
	if entry.Frame != nil {
		*buffer = append(*buffer, ' ')
		if t.FormatCaller != nil {
			t.FormatCaller(buffer, entry.Frame)
		} else {
			FormatCaller(buffer, entry.Frame)
		}
	}
	if len(entry.Message) > 0 {
		*buffer = append(*buffer, ' ')
		*buffer = append(*buffer, entry.Message...)
	}

	if len(entry.Args) > 0 {
		*buffer = append(*buffer, ' ')
		if t.FormatArgs != nil {
			t.FormatArgs(buffer, entry.Args)
		} else {
			*buffer = fmt.Appendln(*buffer, entry.Args...)
		}
	}

	if len(*buffer) == 0 || (*buffer)[len(*buffer)-1] != '\n' {
		*buffer = append(*buffer, '\n')
	}
}
