package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/vmihailenco/msgpack/v5"
)

/* 只複製指定的字串 避免額外的資料注入 */
func CopyParams(source, target map[string]interface{}, keys ...string) {
	for _, key := range keys {
		if source[key] != nil {
			target[key] = source[key]
		}
	}
}

/* 從MAP直接取得頁碼資訊 */
func GetPage(params map[string]interface{}) (pageNo, pageSize int) {
	pageNo = 1
	pageSize = 20
	if val, ok := params["pageNo"].(float64); ok {
		pageNo = int(val)
	}
	if val, ok := params["pageSize"].(float64); ok {
		pageSize = int(val)
	}
	return
}

/* 每個個字串都不為空 回傳true */
func StrNotBlank(str ...string) bool {
	for _, v := range str {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	return true
}

/* 每個字串都為空 回傳true */
func StrAllBlank(str ...string) bool {
	for _, v := range str {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

/* 任何一個字串為空 回傳true */
func StrAnyBlank(str ...string) bool {
	return !StrNotBlank(str...)
}

/* 任何一個字串不為空 回傳true */
func StrAnyNotBlank(str ...string) bool {
	return !StrAllBlank(str...)
}

/* json unmarshal to map */
func JsonToMap(jsonStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{}, 0)
	err := json.Unmarshal([]byte(jsonStr), &m)
	return m, err
}

/* json marshal , skip the step of bytes to string */
func ToJsonStr(obj interface{}) (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

/* Msgpack unmarshal to map with base64*/
func MsgpackToMap(msgpackStr string) (map[string]interface{}, error) {
	m := make(map[string]interface{}, 0)
	b, err := base64.StdEncoding.DecodeString(msgpackStr)
	if err != nil {
		return m, err
	}
	err = msgpack.Unmarshal(b, &m)
	return m, err
}

/* Msgpack marshal with base64 */
func ToMsgpackStr(obj interface{}) (string, error) {
	bytes, err := msgpack.Marshal(obj)
	res := base64.StdEncoding.EncodeToString(bytes)
	return res, err
}

/* convert interface{} to int */
func ToInt(obj interface{}) (int, error) {
	switch d := obj.(type) {
	case float32:
		return int(d), nil
	case float64:
		return int(d), nil
	case int:
		return d, nil
	case int8:
		return int(d), nil
	case int16:
		return int(d), nil
	case int32:
		return int(d), nil
	case int64:
		return int(d), nil
	case uint:
		return int(d), nil
	case uint16:
		return int(d), nil
	case uint32:
		return int(d), nil
	case uint64:
		return int(d), nil
	case string:
		return strconv.Atoi(d)
	default:
		if i, ok := obj.(int); ok {
			return i, nil
		}
	}
	return 0, errors.New("interface{} convert to int fail")
}

/* convert interface{} to int64 */
func ToInt64(obj interface{}) (int64, error) {
	switch d := obj.(type) {
	case float32:
		return int64(d), nil
	case float64:
		return int64(d), nil
	case int:
		return int64(d), nil
	case int8:
		return int64(d), nil
	case int16:
		return int64(d), nil
	case int32:
		return int64(d), nil
	case int64:
		return d, nil
	case uint:
		return int64(d), nil
	case uint16:
		return int64(d), nil
	case uint32:
		return int64(d), nil
	case uint64:
		return int64(d), nil
	case string:
		return strconv.ParseInt(d, 10, 64)
	default:
		if i, ok := obj.(int64); ok {
			return i, nil
		}
	}
	return 0, errors.New("interface{} convert to int64 fail")
}

/* convert interface{} to float64 */
func ToFloat64(obj interface{}) (float64, error) {
	switch d := obj.(type) {
	case float32:
		return float64(d), nil
	case float64:
		return d, nil
	case int:
		return float64(d), nil
	case int8:
		return float64(d), nil
	case int16:
		return float64(d), nil
	case int32:
		return float64(d), nil
	case int64:
		return float64(d), nil
	case uint:
		return float64(d), nil
	case uint16:
		return float64(d), nil
	case uint32:
		return float64(d), nil
	case uint64:
		return float64(d), nil
	case string:
		return strconv.ParseFloat(d, 64)
	default:
		if f, ok := obj.(float64); ok {
			return f, nil
		}
	}
	return 0, errors.New("interface{} convert to float64 fail")
}

/* convert interface{} to str */
func ToStr(obj interface{}) string {
	switch d := obj.(type) {
	case float32:
		return strconv.FormatFloat(float64(d), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(d, 'f', -1, 64)
	case int:
		return strconv.Itoa(d)
	case int8:
		return strconv.FormatInt(int64(d), 10)
	case int16:
		return strconv.FormatInt(int64(d), 10)
	case int32:
		return strconv.FormatInt(int64(d), 10)
	case int64:
		return strconv.FormatInt(d, 10)
	case uint:
		return strconv.FormatUint(uint64(d), 10)
	case uint16:
		return strconv.FormatUint(uint64(d), 10)
	case uint32:
		return strconv.FormatUint(uint64(d), 10)
	case uint64:
		return strconv.FormatUint(d, 10)
	case string:
		return d
	case nil:
		return ""
	default:
		if str, ok := obj.(string); ok {
			return str
		}
		return fmt.Sprint(obj)
	}
}

/* 判段陣列字串裡面是否包含指定字串 */
func ArrayStrInclude(array []string, str string) bool {
	for _, v := range array {
		if v == str {
			return true
		}
	}
	return false
}

/* 組裝query string */
func QueryString(params map[string]interface{}) string {
	values := url.Values{}
	for key, value := range params {
		values.Add(key, ToStr(value))
	}
	return values.Encode()
}

/* queryString to Object */
func QueryStringParse(str string) (map[string]interface{}, error) {
	values, err := url.ParseQuery(str)
	m := make(map[string]interface{})
	for k, v := range values {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m, err
}

/* 錢轉遊戲幣 */
func MoneyToGold(money float64) int64 {
	money *= 100
	s := strconv.FormatFloat(money, 'f', 0, 64)
	resp, _ := strconv.ParseInt(s, 10, 64)
	return resp
}

/* 遊戲幣轉錢 */
func GoldToMoney(money int64) float64 {
	m := float64(money) / 100
	s := strconv.FormatFloat(m, 'f', 2, 64)
	resp, _ := strconv.ParseFloat(s, 64)
	return resp
}
