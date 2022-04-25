package zlog

import (
	"testing"
)

func ExampleInitLog() {
	InitLog("/Users/syss8/OneDrive/Documents/golang/common/zlog/demo_info.log",
		"/Users/syss8/OneDrive/Documents/golang/common/zlog/demo_error.log", false)
	Debug("debug")
	Debugf("%v", "debugf")
	Info("info")
	Infof("%v", "infof")
	Warn("warn")
	Warnf("%v", "warnf")
	Error("error")
	Errorf("%v", "errorf")
	// output:
	//
}

func ExampleConsoleInit() {
	Debug("debug")
	Debugf("%v", "debugf")
	Info("info", "wdasd")
	Infof("%v", "infof")
	Error("error")
	Errorf("%v", "errorf")
	Warn("warn")
	Warnf("%v", "warnf")

	// output:
	//
}

func TestGetSugarLog(t *testing.T) {
	this := GetSugarLog()
	if this == nil {
		t.Error()
	}
}

func TestGetLog(t *testing.T) {
	this := GetLog()
	if this == nil {
		t.Error()
	}
}
