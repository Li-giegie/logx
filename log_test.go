package logx

import (
	"io"
	"testing"
)

func TestLogxText(t *testing.T) {
	l := New()
	l.Debug("Debug", "data1", "date2")
	l.Info("Info", "data2", "date3")
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}

func TestLogxTextPrefixSufxxArgs(t *testing.T) {
	l := New()
	DefaultFormat.(*TextFormat).Color = true
	l.PrefixArgs = append(l.PrefixArgs, "A", "B", "C")
	l.SuffixArgs = append(l.SuffixArgs, "E", "O", "F")
	l.Debug("Debug", "data1", "date2")
	l.Info("Info", "data2", "date3")
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}

func TestLogxTextLevel(t *testing.T) {
	l := New()
	l.Level = LevelInfo
	l.Debug("Debug", "data1", "date2")
	l.Info("Info", "data2", "date3")
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}
func TestLogxHooks(t *testing.T) {
	l := New()
	l.EntryHooks = append(l.EntryHooks, func(l Level, entry *Entry) error {
		// .....
		return nil
	})
	l.BeforeHooks = append(l.BeforeHooks, func(l Level, data []byte) error {
		// ......
		return nil
	})
	l.AfterHooks = append(l.AfterHooks, func(l Level, data []byte) error {
		// ......
		return nil
	})
	l.Debug("Debug", "data1", "date2")
	l.Info("Info", "data2", "date3")
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}
func TestLogxJSON(t *testing.T) {
	l := New()
	l.AddSource = true
	l.Formater = &JSONFormat{
		FormatArgs: FormatJSONArgs,
	}
	l.Debug("Debug", "name", "张三")
	l.Info("Info", "data2", 1)
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}
func TestLogxTextFieldFormat(t *testing.T) {
	l := New()
	l.AddSource = true
	l.Formater = &TextFieldFormat{}
	l.Debug("Debug", "name", "date2")
	l.Info("Info", "data2", "date3")
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}

func BenchmarkLogxTextFormat(b *testing.B) {
	l := New()
	l.Output = io.Discard
	l.AddSource = true
	l.Formater = &TextFormat{}
	for i := 0; i < b.N; i++ {
		l.Info("info", "key", i)
	}
}

func BenchmarkLogxTextFieldFormat(b *testing.B) {
	l := New()
	l.Output = io.Discard
	l.AddSource = true
	l.Formater = &TextFieldFormat{}
	for i := 0; i < b.N; i++ {
		l.Info("info", "key", i)
	}
}

func BenchmarkLogxJSONFormat(b *testing.B) {
	l := New()
	l.Output = io.Discard
	l.AddSource = true
	l.Formater = &JSONFormat{}
	for i := 0; i < b.N; i++ {
		l.Info("info", "key", i)
	}
}
