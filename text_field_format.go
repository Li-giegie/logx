package logx

import (
	"fmt"
	"runtime"
	"time"
)

type TextFieldFormat struct {
	FormatTime   func(buffer *[]byte, t time.Time)
	FormatCaller func(buffer *[]byte, frame *runtime.Frame)
	FormatArgs   func(buffer *[]byte, args []any)
}

func (t *TextFieldFormat) Format(buffer *[]byte, entry *Entry) {
	*buffer = append(*buffer, "time=\""...)
	if t.FormatTime != nil {
		t.FormatTime(buffer, entry.Time)
	} else {
		FormatTime(buffer, entry.Time)
	}
	*buffer = append(*buffer, "\" level=\""...)
	*buffer = append(*buffer, entry.Level.String()...)
	if entry.Frame != nil {
		*buffer = append(*buffer, "\" caller=\""...)
		if t.FormatCaller != nil {
			t.FormatCaller(buffer, entry.Frame)
		} else {
			FormatCaller(buffer, entry.Frame)
		}
	}
	*buffer = append(*buffer, "\" message=\""...)
	if len(entry.Message) > 0 {
		*buffer = append(*buffer, entry.Message...)
	}
	*buffer = append(*buffer, '"')
	if len(entry.Args) > 0 {
		if t.FormatArgs != nil {
			t.FormatArgs(buffer, entry.Args)
		} else {
			for i, arg := range entry.Args {
				if i%2 == 0 {
					*buffer = append(*buffer, ' ')
					*buffer = append(*buffer, arg.(string)...)
					*buffer = append(*buffer, '=')
					continue
				}
				*buffer = fmt.Appendf(*buffer, `"%v"`, arg)
			}
		}
	}
	if len(*buffer) == 0 || (*buffer)[len(*buffer)-1] != '\n' {
		*buffer = append(*buffer, '\n')
	}
}
