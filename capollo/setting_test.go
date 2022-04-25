package capollo

import (
	"testing"
)

func Test_AutoSetting(t *testing.T) {
	Reset()
	set := InitWithConfig(testConfigFilePath)
	tests := []struct {
		name string
	}{
		{"t1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set.AutoSetting()
		})
	}
}
