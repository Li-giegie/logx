package logx

import (
	"io"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	Default       Logger
	DefaultFormat Formater  = new(TextFormat)
	DefaultOutput io.Writer = os.Stdout
	BufferPool              = sync.Pool{
		New: func() any {
			buf := make([]byte, 1024)
			return &buf
		},
	}
	entryPool = sync.Pool{
		New: func() any {
			return new(Entry)
		},
	}
	pcPool = sync.Pool{
		New: func() any {
			pc := make([]uintptr, 1)
			return &pc
		},
	}
)

type Formater interface {
	Format(buffer *[]byte, entry *Entry)
}

func New() *Logger {
	return &Logger{}
}

type Level int

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return "UNKNOWN" + strconv.Itoa(int(l))
	}
}

func (l Level) Color() string {
	switch l {
	case LevelDebug:
		return "\u001B[37mDEBUG\033[0m"
	case LevelInfo:
		return "\u001B[36mINFO\033[0m"
	case LevelWarn:
		return "\u001B[33mWARN\033[0m"
	case LevelError:
		return "\u001B[31mERROR\033[0m"
	default:
		return "UNKNOWN" + strconv.Itoa(int(l))
	}
}

const (
	LevelDebug Level = 0 - iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	Level       Level                               // 日志记录等级
	AddSource   bool                                // 是否记录源码位置
	Output      io.Writer                           // 日志的输出
	EntryHooks  []func(l Level, entry *Entry) error // 日志条目钩子函数，当一条日志被构建完成，还没有格式化时，该切片被一次回调，返回值error不为nil使用println输入错误
	BeforeHooks []func(l Level, data []byte) error  // 日志格式钩子函数，当一条日志被格式化完成，还没有输出前回调
	AfterHooks  []func(l Level, data []byte) error  // 日志格式钩子函数，当一条日志被格式化完成，输出后进行回调
	Formater    Formater                            // 日志的输出格式
	PrefixArgs  []any                               // args第一组参数，设置该参数每次输出日志都会固定输出
	SuffixArgs  []any                               // args最后一组参数，设置该参数每次输出日志都会固定输出
}

func (l *Logger) Debug(msg string, a ...any) {
	if l.Level >= LevelDebug {
		l.Log(LevelDebug, msg, a)
	}
}

func (l *Logger) Info(msg string, a ...any) {
	if l.Level >= LevelInfo {
		l.Log(LevelInfo, msg, a)
	}
}

func (l *Logger) Warn(msg string, a ...any) {
	if l.Level >= LevelWarn {
		l.Log(LevelWarn, msg, a)
	}
}

func (l *Logger) Error(msg string, a ...any) {
	if l.Level >= LevelError {
		l.Log(LevelError, msg, a)
	}
}

// Log 输出日志，level不在判断是否合法
func (l *Logger) Log(level Level, msg string, args []any) {
	entry := entryPool.Get().(*Entry)
	entry.Logger = l
	entry.Level = level
	entry.Time = time.Now()
	entry.Message = msg
	entry.Args = entry.Args[:0]
	entry.Args = append(entry.Args, l.PrefixArgs...)
	entry.Args = append(entry.Args, args...)
	entry.Args = append(entry.Args, l.SuffixArgs...)
	if l.AddSource {
		pc := pcPool.Get().(*[]uintptr)
		runtime.Callers(3, *pc)
		frame, _ := runtime.CallersFrames(*pc).Next()
		pcPool.Put(pc)
		entry.Frame = &frame
	}
	for _, hook := range l.EntryHooks {
		if err := hook(level, entry); err != nil {
			println(time.Now().String(), "entry hook err:", err.Error())
		}
	}
	buffer := BufferPool.Get().(*[]byte)
	*buffer = (*buffer)[:0]
	if l.Formater != nil {
		l.Formater.Format(buffer, entry)
	} else {
		DefaultFormat.Format(buffer, entry)
	}
	for _, hook := range l.BeforeHooks {
		if err := hook(level, *buffer); err != nil {
			println(time.Now().String(), "before hook err:", err.Error())
		}
	}
	if l.Output != nil {
		l.Output.Write(*buffer)
	} else {
		DefaultOutput.Write(*buffer)
	}
	for _, hook := range l.AfterHooks {
		if err := hook(level, *buffer); err != nil {
			println(time.Now().String(), "after hook err:", err.Error())
		}
	}
	entryPool.Put(entry)
	BufferPool.Put(buffer)
}

// Entry 日志条目
type Entry struct {
	*Logger
	Level   Level
	Time    time.Time
	Frame   *runtime.Frame // 调用者帧，只有Logger.AddSource 为true时不为nil
	Message string
	Args    []any
}

func FormatTime(b *[]byte, now time.Time) {
	y, m, d := now.Date()
	h, mm, s := now.Clock()
	*b = append(*b, itoa(y, 4)...)
	*b = append(*b, '-')
	*b = append(*b, itoa(int(m), 2)...)
	*b = append(*b, '-')
	*b = append(*b, itoa(d, 2)...)
	*b = append(*b, 32)
	*b = append(*b, itoa(h, 2)...)
	*b = append(*b, ':')
	*b = append(*b, itoa(mm, 2)...)
	*b = append(*b, ':')
	*b = append(*b, itoa(s, 2)...)
}

func FormatCaller(b *[]byte, frame *runtime.Frame) {
	*b = append(*b, frame.File...)
	*b = append(*b, ':')
	*b = append(*b, strconv.FormatInt(int64(frame.Line), 10)...)
	*b = append(*b, ' ')
	*b = append(*b, frame.Function...)
}

func itoa(i int, wid int) []byte {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	b[bp] = byte('0' + i)
	return b[bp:]
}

func Debug(msg string, a ...any) {
	Default.Log(LevelDebug, msg, a)
}

func Info(msg string, a ...any) {
	Default.Log(LevelInfo, msg, a)
}

func Warn(msg string, a ...any) {
	Default.Log(LevelWarn, msg, a)
}

func Error(msg string, a ...any) {
	Default.Log(LevelError, msg, a)
}
