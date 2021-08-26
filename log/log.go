package log


import (
	"context"
	"fmt"
	context2 "github.com/danakum/util/traceable_context"
	"github.com/google/uuid"
	. "github.com/logrusorgru/aurora"
	"log"
	"os"
	"runtime"
)

var nativeLog *log.Logger
var errorLog *log.Logger

var FileDepth = 2

var (
	fatal = `FATAL`
	err   = `ERROR`
	warn  = `WARN`
	info  = `INFO`
	debug = `DEBUG`
	trace = `TRACE`
)

var logColors = map[string]string{
	`FATAL`: BgRed(`[FATAL]`).String(),
	`ERROR`: BgRed(`[ERROR]`).String(),
	`WARN`:  BgBrown(`[WARN]`).String(),
	`INFO`:  BgBlue(`[INFO]`).String(),
	`DEBUG`: BgCyan(`[DEBUG]`).String(),
	`TRACE`: BgMagenta(`[TRACE]`).String(),
}

var logTypes = map[string]int{
	`FATAL`: 1,
	`ERROR`: 2,
	`WARN`:  3,
	`INFO`:  4,
	`DEBUG`: 5,
	`TRACE`: 6,
}

func init() {
	nativeLog = log.New(os.Stdout, ``, log.LstdFlags|log.Lmicroseconds)
	errorLog = log.New(os.Stderr, ``, log.LstdFlags|log.Lmicroseconds)
}

func colored(typ string) string {
	if Config.Colors {
		return logColors[typ]
	}

	//return `[` + typ + `]`
	return typ
}

//isLoggable Check whether the log type is loggable under current configurations
func isLoggable(logType string) bool {
	return logTypes[logType] <= logTypes[Config.Level]
}

func toString(id string, typ string, message interface{}, params ...interface{}) string {

	var messageFmt = "%s %s %v"

	return fmt.Sprintf(messageFmt,
		typ,
		fmt.Sprintf("%+v", message),
		fmt.Sprintf("%+v", params))
}

func ErrorContext(ctx context.Context, message interface{}, params ...interface{}) {
	logEntryContext(err, ctx, message, colored(`ERROR`), params...)
}

func WarnContext(ctx context.Context, message interface{}, params ...interface{}) {
	logEntryContext(warn, ctx, message, colored(`WARN`), params...)
}

func InfoContext(ctx context.Context, message interface{}, params ...interface{}) {
	logEntryContext(info, ctx, message, colored(`INFO`), params...)
}

func DebugContext(ctx context.Context, message interface{}, params ...interface{}) {
	logEntryContext(debug, ctx, message, colored(`DEBUG`), params...)
}

func TraceContext(ctx context.Context, message interface{}, params ...interface{}) {
	logEntryContext(trace, ctx, message, colored(`TRACE`), params...)
}

func Error(message interface{}, params ...interface{}) {
	logEntry(err, uuid.New(), message, colored(`ERROR`), params...)
}

func Warn(message interface{}, params ...interface{}) {
	logEntry(warn, uuid.New(), message, colored(`WARN`), params...)
}

func Info(message interface{}, params ...interface{}) {
	logEntry(info, uuid.New(), message, colored(`INFO`), params...)
}

func Debug(message interface{}, params ...interface{}) {
	logEntry(debug, uuid.New(), message, colored(`DEBUG`), params...)
}

func Trace(message interface{}, params ...interface{}) {
	logEntry(trace, uuid.New(), message, colored(`TRACE`), params...)
}

func Fatal(message interface{}, params ...interface{}) {
	logEntry(fatal, uuid.New(), message, colored(`FATAL`), params...)
}

func Fataln(message interface{}, params ...interface{}) {
	logEntry(fatal, uuid.New(), message, colored(`FATAL`), params...)
}

func FatalContext(ctx context.Context, message interface{}, params interface{}) {
	logEntry(fatal, uuid.New(), message, colored(`FATAL`), params)
}

func logEntryContext(logType string, ctx context.Context, message interface{}, color string, params ...interface{}) {
	logEntry(logType, uuidFromContext(ctx), message, color, params...)
}

func WithPrefix(p string, message interface{}) string {
	return fmt.Sprintf(`%s] [%+v`, p, message)
}

func uuidFromContext(ctx context.Context) uuid.UUID {
	traceableCtx, ok := ctx.(context2.TraceableContext)
	if !ok {
		return uuid.New()
	}
	return traceableCtx.UUID()
}

func logEntry(logType string, uuid uuid.UUID, message interface{}, color string, params ...interface{}) {

	if !isLoggable(logType) {
		return
	}

	var file string
	var line int
	if Config.FilePath {
		_, f, l, ok := runtime.Caller(FileDepth)
		if !ok {
			f = `<Unknown>`
			l = 1
		}

		file = f
		line = l

		message = fmt.Sprintf(`[%s] [%+v on %s %d]`, uuid.String(), message, file, line)
	} else {
		message = fmt.Sprintf(`[%s] [%+v]`, uuid.String(), message)
	}

	if logType == fatal {
		nativeLog.Fatalln(toString(``, color, message, params...))
	}

	if logType == err {
		nativeLog.Println(toString(``, color, message, params...))
		return
	}

	nativeLog.Println(toString(``, color, message, params...))
}
