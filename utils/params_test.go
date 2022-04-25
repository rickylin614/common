package utils

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestCopyParams(t *testing.T) {
	tests := []struct {
		name       string
		source     map[string]interface{}
		target     map[string]interface{}
		keys       []string
		wantEquals bool
	}{
		{"1", map[string]interface{}{"a": 123, "b": 456}, map[string]interface{}{}, []string{"a", "b"}, true},
		{"2", map[string]interface{}{"a": 123, "b": 456}, map[string]interface{}{}, []string{"a"}, false},
		{"3", map[string]interface{}{"a": 123, "b": 456}, map[string]interface{}{}, []string{"a", "b", "c"}, true},
		{"4", map[string]interface{}{"a": 123, "b": 456}, map[string]interface{}{}, []string{"b", "c"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CopyParams(tt.source, tt.target, tt.keys...)
			if !reflect.DeepEqual(tt.source, tt.target) != !tt.wantEquals {
				t.Errorf("it is error , wantEquals: %v, tt.source: %v, tt.target %v \n", tt.wantEquals, tt.source, tt.target)
			}
		})
	}
}

func TestGetPage(t *testing.T) {
	type args struct {
		params map[string]interface{}
	}
	tests := []struct {
		name         string
		params       map[string]interface{}
		wantPageNo   int
		wantPageSize int
	}{
		{"n1", map[string]interface{}{
			"pageNo":   float64(2),
			"pageSize": float64(40),
		}, 2, 40},
		{"n2", map[string]interface{}{
			"pageNo":   2,
			"pageSize": 40,
		}, 1, 20},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPageNo, gotPageSize := GetPage(tt.params)
			if gotPageNo != tt.wantPageNo {
				t.Errorf("GetPage() gotPageNo = %v, want %v", gotPageNo, tt.wantPageNo)
			}
			if gotPageSize != tt.wantPageSize {
				t.Errorf("GetPage() gotPageSize = %v, want %v", gotPageSize, tt.wantPageSize)
			}
		})
	}
}

func TestStrNotBlank(t *testing.T) {
	tests := []struct {
		name string
		str  []string
		want bool
	}{
		{"n1", []string{"", "", ""}, false},
		{"n1", []string{"1", "2", "3"}, true},
		{"n1", []string{"1", "2", " "}, false},
		{"n1", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrNotBlank(tt.str...); got != tt.want {
				t.Errorf("StrNotBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrAllBlank(t *testing.T) {
	type args struct {
		str []string
	}
	tests := []struct {
		name string
		str  []string
		want bool
	}{
		{"n1", []string{"", "", ""}, true},
		{"n1", []string{"   "}, true},
		{"n1", []string{"1", "2", "3"}, false},
		{"n1", []string{"", "", "", "4"}, false},
		{"n1", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrAllBlank(tt.str...); got != tt.want {
				t.Errorf("StrAllBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrAnyBlank(t *testing.T) {
	tests := []struct {
		name string
		str  []string
		want bool
	}{
		{"n1", []string{"", "", ""}, true},
		{"n1", []string{"1", "2", "3"}, false},
		{"n1", []string{"1", "2", " "}, true},
		{"n1", []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrAnyBlank(tt.str...); got != tt.want {
				t.Errorf("StrAnyBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrAnyNotBlank(t *testing.T) {
	type args struct {
		str []string
	}
	tests := []struct {
		name string
		str  []string
		want bool
	}{
		{"n1", []string{"", "", ""}, false},
		{"n1", []string{"   "}, false},
		{"n1", []string{"1", "2", "3"}, true},
		{"n1", []string{"", "", "", "4"}, true},
		{"n1", []string{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrAnyNotBlank(tt.str...); got != tt.want {
				t.Errorf("StrAnyNotBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonToMap(t *testing.T) {
	tests := []struct {
		name    string
		want    map[string]interface{}
		jsonStr string
		wantErr bool
	}{
		{"n1", map[string]interface{}{
			"test": "123456",
		}, `{"test":"123456"}`, false},
		{"n1", map[string]interface{}{}, `{==========}`, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonToMap(tt.jsonStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("JsonToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JsonToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToJsonStr(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		want    string
		wantErr bool
	}{
		{"n1", map[string]string{
			"test": "123456",
		}, `{"test":"123456"}`, false},
		{"n1", map[string]map[string]int32{
			"level1": {
				"level2": math.MinInt32,
			},
		}, `{"level1":{"level2":-2147483648}}`, false},
		{"n1", map[string]interface{}{
			"test": make(chan int, 0),
		}, ``, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToJsonStr(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJsonStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToJsonStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryString(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		want   string
	}{
		{"t1", map[string]interface{}{"a": 123, "b": 123.5}, "a=123&b=123.5"},
		{"t2", map[string]interface{}{"a": "SSSS", "中文": "DDD"}, `a=SSSS&%E4%B8%AD%E6%96%87=DDD`},
		{"t3", map[string]interface{}{"arr": []string{"A", "B", "C"}}, `arr=%5BA+B+C%5D`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QueryString(tt.params); got != tt.want {
				t.Errorf("QueryString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgpackToMap(t *testing.T) {
	tests := []struct {
		name       string
		msgpackStr string
		want       map[string]interface{}
		wantErr    bool
	}{
		{"t1", "gqZzdHJpbmejYWJjo2ludHs=", map[string]interface{}{"string": "abc", "int": int8(123)}, false},
		{"t2", "123", map[string]interface{}{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MsgpackToMap(tt.msgpackStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("MsgpackToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for k, v := range tt.want {
				if !reflect.DeepEqual(got[k], v) {
					t.Errorf("MsgpackToMap() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestToMsgpackStr(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		want    string
		wantErr bool
	}{
		{"t1", map[string]interface{}{"string": "abc", "int": 123}, "gqZzdHJpbmejYWJjo2ludHs=", false},
		{"t2", make(chan int, 0), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToMsgpackStr(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToMsgpackStr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToMsgpackStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		want    int
		wantErr bool
	}{
		{"t1", nil, 0, true},
		{"t2", float32(1.2), 1, false},
		{"t3", float64(2.4), 2, false},
		{"t4", int8(4), 4, false},
		{"t5", int16(5), 5, false},
		{"t6", int32(6), 6, false},
		{"t7", int64(math.MaxInt64), 9223372036854775807, false},
		{"t8", uint16(1), 1, false},
		{"t9", uint32(math.MaxUint32), 4294967295, false},
		{"t10", uint64(math.MaxUint64), -1, false}, // 超出範圍
		{"t11", "t11", 0, true},
		{"t12", "9223372036854775808", 9223372036854775807, true}, // 超出範圍
		{"t13", "123", 123, false},
		{"t14", 123, 123, false},
		{"t15", uint(123), 123, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		want    int64
		wantErr bool
	}{
		{"t1", nil, 0, true},
		{"t2", float32(1.2), 1, false},
		{"t3", float64(2.4), 2, false},
		{"t4", int8(4), 4, false},
		{"t5", int16(5), 5, false},
		{"t6", int32(6), 6, false},
		{"t7", int64(math.MaxInt64), 9223372036854775807, false},
		{"t8", uint16(1), 1, false},
		{"t9", uint32(math.MaxUint32), 4294967295, false},
		{"t10", uint64(math.MaxUint64), -1, false}, // 超出範圍
		{"t11", "t11", 0, true},
		{"t12", "9223372036854775808", 9223372036854775807, true}, // 超出範圍
		{"t13", "123", 123, false},
		{"t14", 123, 123, false},
		{"t15", uint(123), 123, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToInt64(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name    string
		obj     interface{}
		want    float64
		wantErr bool
	}{
		{"t1", nil, 0, true},
		{"t2", float32(1.2), 1.2000000476837158, false}, // 浮點誤差
		{"t3", float64(2.4), 2.4, false},
		{"t4", int8(4), 4, false},
		{"t5", int16(5), 5, false},
		{"t6", int32(6), 6, false},
		{"t7", int64(math.MaxInt64), 9223372036854775807, false},
		{"t8", uint16(1), 1, false},
		{"t9", uint32(math.MaxUint32), 4294967295, false},
		{"t10", uint64(math.MaxUint64), 1<<64 - 1, false}, // 只紀錄16位數字
		{"t11", "t11", 0, true},
		{"t12", "9223372036854775808", 9223372036854776000, false}, // 只紀錄16位數字
		{"t13", "123", 123, false},
		{"t14", 123, 123, false},
		{"t15", uint(123), 123, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToFloat64(tt.obj)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToFloat64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToFloat64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStr(t *testing.T) {
	tests := []struct {
		name string
		obj  interface{}
		want string
	}{
		{"t1", nil, ""},
		{"t2", float32(1.2), "1.2000000476837158"}, // 浮點誤差
		{"t3", float64(2.4), "2.4"},
		{"t4", int8(4), "4"},
		{"t5", int16(5), "5"},
		{"t6", int32(6), "6"},
		{"t7", int64(math.MaxInt64), "9223372036854775807"},
		{"t8", uint16(1), "1"},
		{"t9", uint32(math.MaxUint32), "4294967295"},
		{"t10", uint64(math.MaxUint64), "18446744073709551615"},
		{"t11", "9223372036854775808", "9223372036854775808"},
		{"t12", 123, "123"},
		{"t13", uint(123), "123"},
		{"t14", map[string]string{}, "map[]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToStr(tt.obj); got != tt.want {
				t.Errorf("ToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayStrInclude(t *testing.T) {
	tests := []struct {
		name  string
		array []string
		str   string
		want  bool
	}{
		{"t1", []string{"a"}, "b", false},
		{"t2", []string{"a", "b"}, "b", true},
		{"t3", []string{"b"}, "b", true},
		{"t4", []string{"b", "c"}, "b", true},
		{"t5", []string{"a", "b", "c"}, "d", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ArrayStrInclude(tt.array, tt.str); got != tt.want {
				t.Errorf("ArrayStrInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryStringParse(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		want    map[string]interface{}
		wantErr bool
	}{
		{"n1", "a=123&b=123.5", map[string]interface{}{"a": 123, "b": 123.5}, false},
		{"n2", `a=SSSS&%E4%B8%AD%E6%96%87=DDD`, map[string]interface{}{"a": "SSSS", "中文": "DDD"}, false},
		{"n3", `%99%A0%GF`, map[string]interface{}{}, true},
		{"n4", `arr=%5BA+B+C%5D`, map[string]interface{}{"arr": []string{"A", "B", "C"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryStringParse(tt.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryStringParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Printf("%v%T\n", got, got)
			if !reflect.DeepEqual(ToStr(got), ToStr(tt.want)) {
				t.Errorf("QueryStringParse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMoneyToGold(t *testing.T) {
	tests := []struct {
		name  string
		money float64
		want  int64
	}{
		{"t1", 9.99, 999},
		{"t2", -0.01, -1},
		{"t3", -2, -200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MoneyToGold(tt.money); got != tt.want {
				t.Errorf("MoneyToGold() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGoldToMoney(t *testing.T) {
	tests := []struct {
		name  string
		money int64
		want  float64
	}{
		{"t1", 999, 9.99},
		{"t2", -1, -0.01},
		{"t3", -200, -2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GoldToMoney(tt.money); got != tt.want {
				t.Errorf("GoldToMoney() = %v, want %v", got, tt.want)
			}
		})
	}
}
