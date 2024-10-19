package logger

type Logger interface {
	Info(wrappedMsg string, fields ...any)
	Debug(wrappedMsg string, fields ...any)
	Error(wrappedMsg string, fields ...any)
	Fatal(wrappedMsg string, fields ...any)
}
