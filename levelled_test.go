package levelled_test

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/krocos/levelled"
)

func ExampleLogger_Error() {
	// Here we wrap zap logger to provide severity level (here it is the error lvl)
	llogger := levelled.New(zap.NewExample(), zapcore.ErrorLevel)

	// These logs will not be logged until error log occur
	llogger.Debug("first")
	llogger.Info("second")
	llogger.Debug("third")

	// Also, you can erase previously written logs
	// to use levelled logger, for example, in cycles
	//llogger.Erase()

	// Here we sleep a bit to show that time of each entry
	// is preserved (if you turn on timestamp writing)
	time.Sleep(time.Second)

	// Create separate levelled logger with additional fields
	// to show that it preserves previous not logged entries
	llogger = llogger.With(zap.String("additional", "field"))

	llogger.Warn("fourth")
	llogger.Info("fifth")
	llogger.Debug("sixth")

	time.Sleep(time.Second)

	// When error had occurred so that all
	// previous logs will be written to their sinks
	llogger.Error("error")

	// Output:
	// {"level":"debug","msg":"first"}
	// {"level":"info","msg":"second"}
	// {"level":"debug","msg":"third"}
	// {"level":"warn","msg":"fourth","additional":"field"}
	// {"level":"info","msg":"fifth","additional":"field"}
	// {"level":"debug","msg":"sixth","additional":"field"}
	// {"level":"error","msg":"error","additional":"field"}
}
