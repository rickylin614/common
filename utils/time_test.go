package utils

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/rickylin614/common/constants"
	"github.com/rickylin614/common/zlog"
)

func ExampleGetTimeString() {
	t := time.Now()
	defer TimeTrack(t)
	time.Sleep(time.Millisecond)
	s1 := GetTimeString(time.Now(), "Europe/Malta", constants.TimeFormat.FLOAT03())   // GMT+2
	s2 := GetTimeString(time.Now(), "Asia/Hong_Kong", constants.TimeFormat.FLOAT03()) // GMT+8
	s3 := GetTimeString(time.Now(), "Asia/Tokyo", constants.TimeFormat.FLOAT03())     // GMT+9
	s4 := GetTimeString(time.Now(), "AAA", constants.TimeFormat.FLOAT03())            //error
	zlog.Info(s1)
	zlog.Info(s2)
	zlog.Info(s3)
	zlog.Info(s4)
	zlog.Info(GetMilli(t))
	zlog.Info(GetMilliStr(t))
	zlog.Info(GetTimeStr(t))

	// output:
	//
}

func TestSqlTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		t    []time.Time
		want sql.NullTime
	}{
		{"n1", []time.Time{now}, sql.NullTime{Time: now, Valid: true}},
		{"n2", nil, sql.NullTime{Valid: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SqlTime(tt.t...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SqlTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayStartTime(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{
		{"t1", time.Unix(0, 1630987508*int64(time.Second)), time.Unix(0, 1630965600*int64(time.Second))},
		{"t2", time.Unix(0, 1612367984*int64(time.Second)), time.Unix(0, 1612306800*int64(time.Second))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DayStartTime(tt.t); got != tt.want {
				t.Errorf("DayStartTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateStr(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want string
	}{
		{"t1", time.Unix(0, 1630987508*int64(time.Second)), "20210907"},
		{"t2", time.Unix(0, 987508*int64(time.Second)), "19700112"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DateStr(tt.t); got != tt.want {
				t.Errorf("DateStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsToTime(t *testing.T) {
	tests := []struct {
		name string
		ms   int64
		want time.Time
	}{
		{"t1", 1630987508000, time.Unix(1630987508, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MsToTime(tt.ms); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNsToTime(t *testing.T) {
	tests := []struct {
		name string
		ns   int64
		want time.Time
	}{
		{"t1", 1630987508000000000, time.Unix(1630987508, 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NsToTime(tt.ns); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NsToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
