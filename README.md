# Logx

一个简单、易扩展，高性能的Go日志框架。

## 特征
- 简单
- 高性能、超小内存分配
- 支持等级
- 支持Hook
- 支持JSON、Text、和其他的扩展方式

## 导入
import
```
import "github.com/Li-giegie/logx"
```
go get
```
go get github.com/Li-giegie/logx
```

## 用法


`注意默认日志等级为Debug`
```go
type Logger struct {
	Level       Level                               // 日志记录等级
	AddSource   bool                                // 是否记录源码位置
	Output      io.Writer                           // 日志的输出
	EntryHooks  []func(l Level, entry *Entry) error // 日志条目钩子函数，当一条日志被构建完成，还没有格式化时，该切片被一次回调，返回值error不为nil使用println输入错误
	BeforeHooks []func(l Level, data []byte) error  // 日志格式钩子函数，当一条日志被格式化完成，还没有输出前回调
	AfterHooks  []func(l Level, data []byte) error  // 日志格式钩子函数，当一条日志被格式化完成，输出后进行回调
	Formater    Formater                            // 日志的输出格式
}

func TestLogxText(t *testing.T) {
	l := New()
	l.Debug("Debug", "data1", "date2")
	l.Info("Info", "data2", "date3")
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}
```