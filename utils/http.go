package utils

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rickylin614/common/zlog"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmhttp"
	"golang.org/x/net/context/ctxhttp"
)

var HttpMockMap map[string]map[string]interface{} = make(map[string]map[string]interface{})
var HttpMockErr map[string]error = make(map[string]error)

var tracingClient = apmhttp.WrapClient(http.DefaultClient)

func HttpGetJSON(ctx context.Context, url string,
	param map[string]interface{}) (resMap map[string]interface{}, err error) {

	//Mock用
	if HttpMockMap[url] != nil {
		return HttpMockMap[url], nil
	}
	if HttpMockErr[url] != nil {
		return nil, HttpMockErr[url]
	}

	//實際執行 組合查詢條件參數
	if qs := QueryString(param); qs != "" {
		url = url + "?" + qs
	}

	//打GET請求出去
	resp, err := ctxhttp.Get(ctx, tracingClient, url)
	if err != nil {
		apm.CaptureError(ctx, err).Send()
		zlog.Errorf("http Client error, url:%s , err:%v\n", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	//讀取回傳值
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//判斷回傳200 OK
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}
	//解析JSON 綁在MAP內回傳
	err = json.Unmarshal(body, &resMap)
	return resMap, err
}
