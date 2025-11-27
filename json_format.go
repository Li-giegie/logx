package logx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

type JSONFormat struct {
	FormatTime   func(buffer *[]byte, t time.Time)
	FormatCaller func(buffer *[]byte, frame *runtime.Frame)
	FormatArgs   func(buffer *[]byte, args []any)
}

func (f *JSONFormat) Format(buffer *[]byte, entry *Entry) {
	*buffer = append(*buffer, `{"level":"`...)
	*buffer = append(*buffer, entry.Level.String()...)
	*buffer = append(*buffer, `","time":"`...)
	if f.FormatTime != nil {
		f.FormatTime(buffer, entry.Time)
	} else {
		*buffer = entry.Time.AppendFormat(*buffer, time.RFC3339Nano)
	}
	*buffer = append(*buffer, `","caller":"`...)
	if entry.Logger.AddSource {
		if f.FormatCaller != nil {
			f.FormatCaller(buffer, entry.Frame)
		} else {
			FormatCaller(buffer, entry.Frame)
		}
	}
	*buffer = append(*buffer, `","message":"`...)
	*buffer = append(*buffer, entry.Message...)
	*buffer = append(*buffer, `","args":`...)
	if f.FormatArgs != nil {
		f.FormatArgs(buffer, entry.Args)
	} else {
		writer := bytes.NewBuffer(*buffer)
		err := json.NewEncoder(writer).Encode(entry.Args)
		if err != nil {
			println(time.Now().String(), "json format encode args err", err.Error())
		}
		if writer.Bytes()[writer.Len()-1] == '\n' {
			*buffer = writer.Bytes()[:writer.Len()-1]
		} else {
			*buffer = writer.Bytes()
		}
	}
	*buffer = append(*buffer, "}\n"...)
}

func FormatJSONArgs(buffer *[]byte, args []any) {
	*buffer = append(*buffer, '{')
	for i := 0; len(args) > 0; i++ {
		if i == 0 {
			*buffer = fmt.Appendf(*buffer, `"%v":`, args[0])
		} else {
			*buffer = fmt.Appendf(*buffer, `,"%v":`, args[0])
		}
		args = args[1:]
		if len(args) == 0 {
			*buffer = append(*buffer, "null}"...)
			return
		} else {
			writer := bytes.NewBuffer(*buffer)
			err := json.NewEncoder(writer).Encode(args[0])
			if err != nil {
				println(time.Now().String(), "json encode args err", err.Error())
			}
			if writer.Bytes()[writer.Len()-1] == '\n' {
				*buffer = writer.Bytes()[:writer.Len()-1]
			} else {
				*buffer = writer.Bytes()
			}
			args = args[1:]
		}
	}
	*buffer = append(*buffer, "}"...)
}
