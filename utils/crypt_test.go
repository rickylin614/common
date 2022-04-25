package utils

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"
)

func TestAesEncryptCBC(t *testing.T) {
	m := map[string]interface{}{
		"channel":     0,
		"agent":       "00001",
		"accountId":   "987654321",
		"accountType": 0,
		"moneyType":   0,
		"money":       12345,
		"orderId":     "wtf?",
		"lineCode":    "cccccccc",
	}
	testByteArr, _ := json.Marshal(m)

	tests := []struct {
		name          string
		origData      []byte
		key           []byte
		wantEncrypted []byte
		wantErr       bool
	}{
		{"name", []byte("中文字"), []byte("c42d9026da8a10b834b28bf36db3a6df"), []byte{101, 60, 115, 222, 0, 129, 247, 54, 143, 106, 187, 5, 205, 107, 70, 66}, false},
		{"name", []byte("中文字"), []byte("c42d9026da8a10b834b28bf36db3a6df2"), []byte{101, 60, 115, 222, 0, 129, 247, 54, 143, 106, 187, 5, 205, 107, 70, 66}, true},
		{"name", testByteArr, []byte("c42d9026da8a10b834b28bf36db3a6df"), []byte{206, 26, 27, 123, 250, 16, 197, 90, 255, 237, 166, 52, 232, 17, 9,
			193, 4, 213, 13, 167, 192, 117, 64, 236, 28, 106, 193, 114, 0, 84, 81, 141, 161, 149, 2, 17, 199, 16, 247, 250, 171, 68, 61, 227, 110, 33, 218, 15, 148, 226,
			35, 61, 186, 180, 3, 1, 107, 43, 5, 163, 15, 251, 244, 218, 201, 210, 127, 203, 252, 34, 242, 118, 22, 186, 41, 166, 81, 97, 11, 198, 5, 218, 245, 168, 141,
			122, 185, 55, 223, 79, 13, 93, 4, 190, 14, 18, 64, 107, 226, 68, 45, 113, 175, 51, 18, 28, 210, 126, 181, 2, 100, 230, 246, 114, 120, 165, 110, 31, 42, 224,
			102, 48, 116, 149, 252, 142, 242, 236, 136, 200, 246, 65, 13, 60, 159, 163, 113, 217, 36, 159, 145, 146, 49, 35}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEncrypted, err := AesEncryptCBC(tt.origData, tt.key); !reflect.DeepEqual(gotEncrypted, tt.wantEncrypted) {
				if err != nil {
					if !tt.wantErr {
						t.Errorf("%v", err)
					}
					return
				}
				t.Errorf("AesEncryptCBC() = %v, want %v", gotEncrypted, tt.wantEncrypted)
			}
		})
	}
}

func TestAesDecryptCBC(t *testing.T) {
	type args struct {
		encrypted []byte
		key       []byte
	}
	tests := []struct {
		name          string
		encrypted     []byte
		key           []byte
		wantDecrypted []byte
		wantErr       bool
	}{
		{"t1", []byte{101, 60, 115, 222, 0, 129, 247, 54, 143, 106, 187, 5, 205, 107, 70, 66}, []byte("c42d9026da8a10b834b28bf36db3a6df"), []byte("中文字"), false},
		{"t2", []byte{101, 60, 115, 222, 0, 129, 247, 54, 143, 106, 187, 5, 205, 107, 70, 66}, []byte("c42d9026da8a10b834b28bf36db3a6df2"), []byte("中文字"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDecrypted, err := AesDecryptCBC(tt.encrypted, tt.key); !reflect.DeepEqual(gotDecrypted, tt.wantDecrypted) {
				if err != nil {
					if !tt.wantErr {
						t.Errorf("%v", err)
					}
					return
				}
				t.Errorf("AesDecryptCBC() = %s, want %s", gotDecrypted, tt.wantDecrypted)
			}
		})
	}
}

func TestAesEncryptECB(t *testing.T) {
	tests := []struct {
		name          string
		origData      string
		key           string
		wantEncrypted string
		wantErr       bool
	}{
		{"t1", `{"aa":123,"bb":124}`, "123456789ABCDEFG", "arfJv3uBQNSubOEYFxLbvkh1sCGevYTb/129Qj5RZCw=", false},
		{"t2", `{"aa":アイウエオ,"bb":中文字}`, "123456789ABCDEFG", "/Xka7JJJY6HlXvYmBnINLHelSsU5PMQTA+EvcSkk1c9ZJJ9GFsmQAJKbXyABcZDL", false},
		{"t3", `{"aa":アイウエオ,"bb":中文字}`, "", "", true},
		{"t4", `{"aa":アイウエオ,"bb":中文字}`, "987654321ZYXW331231", "", true},
		{"t5", `{"aa":アイウエオ,"bb":中文字}`, "123456789ABCDEFG123456789ABCDEFG", "XTnBM9kj1D/pKyqs/eT7cKz/ENo6Iv0rzXBz+s3ZZ8DvtBsO9i9uCPvDKgXORy7A", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEncrypted, err := AesEncryptECB([]byte(tt.origData), []byte(tt.key))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v", err)
				}
				return
			}
			got := base64.StdEncoding.EncodeToString(gotEncrypted)
			if tt.wantEncrypted != got {
				t.Errorf("AesEncryptECB() = %v, want %v", got, tt.wantEncrypted)
			}
		})
	}
}

func TestAesDecryptECB(t *testing.T) {
	type args struct {
		encrypted []byte
		key       []byte
	}
	tests := []struct {
		name          string
		wantDecrypted string
		key           string
		encrypted     string
		wantErr       bool
	}{
		{"t1", `{"aa":123,"bb":124}`, "123456789ABCDEFG", "arfJv3uBQNSubOEYFxLbvkh1sCGevYTb/129Qj5RZCw=", false},
		{"t2", `{"aa":アイウエオ,"bb":中文字}`, "123456789ABCDEFG", "/Xka7JJJY6HlXvYmBnINLHelSsU5PMQTA+EvcSkk1c9ZJJ9GFsmQAJKbXyABcZDL", false},
		{"t3", ``, "", "", true},
		{"t4", ``, "987654321ZYXW331231", "", true},
		{"t5", `{"aa":アイウエオ,"bb":中文字}`, "123456789ABCDEFG123456789ABCDEFG", "XTnBM9kj1D/pKyqs/eT7cKz/ENo6Iv0rzXBz+s3ZZ8DvtBsO9i9uCPvDKgXORy7A", false},
		{"t6", ``, "123456789ABCDEFG123456789ABCDEFG", "arfJv3uBQNSubOEYFxLbvkh1sCGevYTb/129Qj5RZCw=", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := base64.StdEncoding.DecodeString(tt.encrypted)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v", err)
				}
				return
			}
			dec, err := AesDecryptECB(b, []byte(tt.key))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v", err)
				}
				return
			}
			if tt.wantDecrypted != string(dec) {
				t.Errorf("AesDecryptECB() = %s, want %v", dec, tt.wantDecrypted)
			}
		})
	}
}

func TestAesEncryptCFB(t *testing.T) {
	tests := []struct {
		name     string
		origData string
		key      string
		wantErr  bool
	}{
		{"t1", `{"aa":123,"bb":124}`, "123456789ABCDEFG", false},
		{"t2", `{"aa":アイウエオ,"bb":中文字}`, "123456789ABCDEFG", false},
		{"t3", ``, "", true},
		{"t4", `{"aa":アイウエオ,"bb":中文字}`, "987654321ZYXW331231", true},
		{"t5", `1`, "123456789ABCDEFG123456789ABCDEFG", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 每次產生的都不同 因此不做確認
			_, err := AesEncryptCFB([]byte(tt.origData), []byte(tt.key))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v", err)
				}
				return
			}
		})
	}
}

func TestAesDecryptCFB(t *testing.T) {
	type args struct {
		encrypted []byte
		key       []byte
	}
	tests := []struct {
		name          string
		wantDecrypted string
		key           string
		encrypted     string
		wantErr       bool
	}{
		{"t1", `{"aa":123,"bb":124}`, "123456789ABCDEFG", "XwKdxDGR0HJsyMOhie/ZxptVhEk8phlr7vRAHVaDjwDttZA=", false},
		{"t2", `{"aa":123,"bb":124}`, "123456789ABCDEFG", "aBjDnh4/WjSfdcDeSbflZ3HVZ+dMgYLX01Vj7cWEM/9/N5M=", false},
		{"t3", `{"aa":123,"bb":124}`, "123456789ABCDEFG", "ePy5yHH0mhQgBZoFO40OZl8MEOcyJSP4Xa8dklJ3cElkqNw=", false},
		{"t4", `{"aa":123,"bb":124}`, "123456789ABCDEFG123456789ABCDEFG", "YWJj", true},
		{"t5", `{"aa":アイウエオ,"bb":中文字}`, "", "", true},
		{"t6", `{"aa":アイウエオ,"bb":中文字}`, "987654321ZYXW331231", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := base64.StdEncoding.DecodeString(tt.encrypted)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v", err)
				}
				return
			}
			dec, err := AesDecryptCFB(b, []byte(tt.key))
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%v", err)
				}
				return
			}
			if tt.wantDecrypted != string(dec) {
				t.Errorf("AesDecryptECB() = %s, want %v", dec, tt.wantDecrypted)
			}
		})
	}
}

func TestMd5Encrypt(t *testing.T) {
	tests := []struct {
		name     string
		origData []string
		want     string
	}{
		{"t1", []string{"12345", "9*999"}, "43d13e394c3aa9f78c78449272c8f69c"},
		{"t1", []string{}, "d41d8cd98f00b204e9800998ecf8427e"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Md5Encrypt(tt.origData...); got != tt.want {
				t.Errorf("Md5Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesBase64Encrypt(t *testing.T) {
	tests := []struct {
		name     string
		origData []byte
		key      []byte
		aesFunc  aesCrypt
		want     string
		wantErr  bool
	}{
		{"t1", []byte("YOYO"), []byte("123456789ABCDEFG"), AesEncryptCBC, "sCIdlnomB/xZVI1Ll/tMfg==", false},
		{"t2", []byte("YOYO"), []byte("123456789ABCDEFG"), AesEncryptECB, "BK2xtB+MXsPuQyDXmS5W7A==", false},
		// {"t3", []byte("YOYO"), []byte("123456789ABCDEFG"), AesEncryptCFB, "e6UCDdGDfdXZr+BYkqD9ShwjLYA=", false}, //CFB 每次結
		{"t4", []byte("YOYO"), []byte("123456789ABCDEF"), AesEncryptCFB, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AesBase64Encrypt(tt.origData, tt.key, tt.aesFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("AesBase64Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AesBase64Encrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAesBase64Decrypt(t *testing.T) {
	tests := []struct {
		name     string
		want     []byte
		key      []byte
		aesFunc  aesCrypt
		origData string
		wantErr  bool
	}{
		{"t1", []byte("YOYO"), []byte("123456789ABCDEFG"), AesDecryptCBC, "sCIdlnomB/xZVI1Ll/tMfg==", false},
		{"t2", []byte("YOYO"), []byte("123456789ABCDEFG"), AesDecryptECB, "BK2xtB+MXsPuQyDXmS5W7A==", false},
		{"t3", []byte("YOYO"), []byte("123456789ABCDEFG"), AesDecryptCFB, "e6UCDdGDfdXZr+BYkqD9ShwjLYA=", false},
		{"t4", nil, []byte("123456789A"), AesDecryptCBC, "AA", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AesBase64Decrypt(tt.origData, tt.key, tt.aesFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("AesBase64Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AesBase64Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
