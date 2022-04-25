package capollo

import (
	"fmt"
	"reflect"
	"testing"
)

const testConfigFilePath = "./test/config.yml"

/* for unit test clear setting */
func Reset() {
	cli = nil
}

func ExampleInitWithConfig() {
	Reset()
	InitWithConfig(testConfigFilePath)
	// output:
	// 初始化apollo設定 appid=go,host=http://10.1.1.152:28080,namespace=GOLANG.api,cluster=DEV,secretkey=40455bf8e6e749178421e58808a7f490
}

func TestInitWithConfig(t *testing.T) {
	tests := []struct {
		name    string
		ymlFile string
	}{
		{"test1", ""},
		{"test1", testConfigFilePath},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitWithConfig(tt.ymlFile)
		})
	}
}

func TestInitApollo(t *testing.T) {
	Reset()
	type args struct {
		appid     string
		host      string
		namespace string
		cluster   string
		secretkey string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "1", args: args{appid: "kyllc", host: "http://localhost:8080", namespace: "application", secretkey: "312f6871d8ff456d834dddcde3a74dd7"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitApollo(tt.args.appid, tt.args.host, tt.args.namespace, tt.args.cluster, tt.args.secretkey)
		})
	}
}

func TestGetValue(t *testing.T) {
	Reset()
	// 判斷沒設定時，能正確地抱錯
	t.Run("test", func(t *testing.T) {
		_, err := GetValue("")
		if (err != nil) != true {
			t.Errorf("GetYmlValue() error = %v, wantErr %v", err, true)
			return
		}
	})

	InitWithConfig(testConfigFilePath)
	tests := []struct {
		value string
		want  string
	}{
		{value: "mysql.yml", want: "host: 10.1.1.152:3306\nuser: root\npwd: 123456"},
		{value: "noSetting", want: ""},
	}
	for _, tt := range tests {
		t.Run("run", func(t *testing.T) {
			str, err := GetValue(tt.value)
			if err != nil {
				t.Error(err)
			}
			if str != tt.want {
				t.Errorf("GetValue() value = %v, want = %v", str, tt.want)
			}
		})
	}
}

func TestGetYmlValue(t *testing.T) {
	Reset()
	// 判斷沒設定時，能正確地抱錯
	t.Run("test", func(t *testing.T) {
		_, err := GetYmlValue("")
		if (err != nil) != true {
			t.Errorf("GetYmlValue() error = %v, wantErr %v", err, true)
			return
		}
	})

	InitWithConfig(testConfigFilePath)
	tests := []struct {
		name    string
		value   string
		want    map[string]interface{}
		wantErr bool
	}{
		{value: "gameConfig", wantErr: false, want: map[string]interface{}{
			"host": "10.1.1.152",
			"port": "3306",
			"user": "root",
			"pwd":  "123456",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetYmlValue(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetYmlValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			g := fmt.Sprint(got)
			want := fmt.Sprint(tt.want)
			if !reflect.DeepEqual(g, want) {
				t.Errorf("GetYmlValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
