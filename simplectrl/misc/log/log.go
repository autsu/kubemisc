package main

import (
	"context"
	"flag"
	stdlog "log"
	"os"

	"github.com/go-logr/stdr"
	"go.uber.org/zap/zapcore"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	SetLogger = log.SetLogger
	Log       = log.FromContext(context.Background())
	funcNum   = flag.Int("i", 0, "")
)

func Zap() {
	// 注意这个 zap 不是 uber 提供的 zap，而是 controller-runtime 下的 zap，
	// 其实就是提供了一层封装，使得其能够适配 logr
	// logr 貌似是 controller-runtime 使用的日志前端，需要调用 log.SetLogger
	// 来设置对应的日志后端
	//
	// 封装的 zap 提供了 UseDevMode 来设置日志配置，类似 zap.NewDevelopment
	SetLogger(zap.New(zap.UseDevMode(false)))
	Log.Info("test", "123", 123)
	// 如果为 Error 日志会输出相应的堆栈
	Log.Error(nil, "123", "123", 123)

	// 这里不会被覆盖，也就是说不会用 stdlog 作为新的 logger，依然沿用上面的 zap
	SetLogger(stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)))
	Log.Info("test", "123", 123)
	Log.Error(nil, "123", "123", 123)
}

func ZapOption() {
	logger := zap.New(zap.UseFlagOptions(&zap.Options{
		Development:     true,
		StacktraceLevel: zapcore.InfoLevel | zapcore.DebugLevel | zapcore.WarnLevel | zapcore.ErrorLevel,
	}))
	SetLogger(logger)

	Log.Info("test", "123", 123)
	Log.Error(nil, "123", "123", 123)
}

func Std() {
	SetLogger(stdr.New(stdlog.New(os.Stderr, "", stdlog.LstdFlags|stdlog.Lshortfile)))
	Log.Info("test", "123", 123)
	Log.Error(nil, "123", "123", 123)
}

func main() {
	flag.Parse()

	switch *funcNum {
	case 1:
		Zap()
	case 2:
		ZapOption()
	case 3:
		Std()
	default:
		panic("invalid input")
	}
}
