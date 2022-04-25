package utils

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/rickylin614/common/constants"
	"github.com/rickylin614/common/zlog"
)

/* 配合defer, 測試時顯示執行時間使用 */
func TimeTrack(start time.Time) {
	elapsed := time.Since(start)
	zlog.Info("took %s\n", elapsed)
}

/* 取得時間 毫秒int64 */
func GetMilli(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

/* 取得時間 毫秒字串 */
func GetMilliStr(t time.Time) string {
	return strconv.FormatInt(t.UnixNano()/int64(time.Millisecond), 10)
}

/* 取得時間字串 yyyyMMddHHmmssSSS */
func GetTimeStr(t time.Time) string {
	s := t.Format(constants.TimeFormat.SYMBOL03())
	s = strings.Replace(s, ".", "", 1)
	return s
}

/* 取得時間字串(需要附帶時區) */
func GetTimeString(t time.Time, location, format string) (str string) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		zlog.Debug("time location error:", err)
		str = t.Format(format)
	} else {
		str = t.In(loc).Format(format)
	}
	return
}

/* 便捷組出sql.NullTime */
func SqlTime(t ...time.Time) sql.NullTime {
	if len(t) > 0 {
		return sql.NullTime{
			Time:  t[0],
			Valid: true,
		}
	} else {
		return sql.NullTime{
			Valid: false,
		}
	}
}

/*
	取得時間當天起始值
*/
func DayStartTime(t time.Time) time.Time {
	nt := time.Date(t.Year(), time.Month(t.Month()), t.Day(), 0, 0, 0, 0, t.Location())
	return nt
}

/*
	取得時間當天日字串
*/
func DateStr(t time.Time) string {
	return t.Format("20060102")
}

/*
	Milli seconds to time
*/
func MsToTime(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}

/*
	Nano secounds to time
*/
func NsToTime(ns int64) time.Time {
	return time.Unix(0, ns)
}
