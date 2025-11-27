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
func TestLogxText(t *testing.T) {
	l := New()
	l.Debug("Debug", "data1", "date2")
	l.Info("Info", "data2", "date3")
	l.Warn("Warn", "data3", "date4")
	l.Error("Error", "data4", "date5")
}
```