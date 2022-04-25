package zlog

import (
	"io"
	"strings"
	"time"

	. "github.com/rickylin614/common/constants"
	"go.elastic.co/apm/module/apmzap"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugarLogger *zap.SugaredLogger
var logger *zap.Logger

/*
	infoFileName: info檔案路徑
	errorFilePath: error檔案路徑
	Levels:可自定義更改預設的infoLevel/ErrorLevel等級
*/
func InitLog(infoFilePath, errorFilePath string, isJson bool, levels ...zapcore.Level) {
	// 設定一些基本日誌格式
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(TimeFormat.FLOAT3()))
		},
		CallerKey:     "caller",
		StacktraceKey: "trace",
		EncodeCaller:  zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	// 設置使用json / console格式輸出
	var encoder zapcore.Encoder
	if isJson {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	level1 := zapcore.InfoLevel
	level2 := zapcore.ErrorLevel

	// 可自定義更改預設的infoLevel/ErrorLevel等級
	if len(levels) > 0 {
		level1 = levels[0]
	}
	if len(levels) > 1 {
		level2 = levels[1]
	}

	// 實現兩個判斷日誌等級的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level1 && lvl < level2
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= level2
	})

	// 獲取 info、error日誌檔案的io.Writer 抽象 getWriter() 在下方實現
	infoWriter := getWriter(infoFilePath)
	errorWriter := getWriter(errorFilePath)

	// 最後建立具體的Logger
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
	)

	log := zap.New(core,
		zap.AddStacktrace(zapcore.WarnLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.WrapCore((&apmzap.Core{}).WrapCore),
	)
	logger = log
	sugarLogger = log.Sugar()
}

func ConsoleInit() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(TimeFormat.FLOAT03()))
	}
	log, err := config.Build()
	if err != nil {
		Panic(err)
	}
	logger = log
	sugarLogger = log.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

// %A 星期名全稱(原生的)
// %a 3個字元的星期名(原生的)
// %B 月份名的全稱(原生的)
// %b 3個字元的月份名(原生的)
// %c 日期和時間(原生的)
// %d 2位數的一個月中的日期數
// %H 2位數的小時數(24小時制)
// %I 2位數的小時數(12小時制)
// %j 3位數的一年中的日期數
// %M 2位數的分鐘數
// %m 2位數的月份數
// %p am/pm12小時制的上下午(原生的)
// %S 2位數的秒數
// %U 2位數的一年中的星期數(星期天爲一週的第一天)
// %W 2位數的一年中的星期數(星期一爲一週的第一天)
// %w 1位數的星期幾(星期天爲一週的第一天)
// %X 時間(原生的)
// %x 日期(原生的)
// %Y 4位元數的年份
// %y 2位數的年份
// %Z 時區名
// %% 符號"%"本身
func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 實際生成的檔名 demo_info-YYYYmmddHH.log
	// demo.log是指向最新日誌的連結
	hook, err := rotatelogs.New(
		strings.Replace(filename, ".log", "", -1)+"-%Y-%m-%d.log", // format格式
		rotatelogs.WithLinkName(filename),                         // 生成連結，指向最新的LOG
		rotatelogs.WithMaxAge(time.Hour*24*7),                     // 存活時間，保存7日
		rotatelogs.WithRotationTime(time.Hour),                    // 日誌切割時間的時間間格
	)

	if err != nil {
		panic(err)
	}
	return hook
}
func Debug(args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Panicf(template, args...)
}
func Fatal(args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	if sugarLogger == nil {
		ConsoleInit()
	}
	sugarLogger.Fatalf(template, args...)
}

func GetSugarLog() *zap.SugaredLogger {
	if sugarLogger == nil {
		ConsoleInit()
	}
	return sugarLogger
}

func GetLog() *zap.Logger {
	if sugarLogger == nil {
		ConsoleInit()
	}
	return logger
}
