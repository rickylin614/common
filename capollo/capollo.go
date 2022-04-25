package capollo

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/rickylin614/common/utils"
	"github.com/rickylin614/common/zlog"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var (
	// 保存連線設定
	cli *agollo.Client
	// 保存namespace
	Namespace *string
	// Unit Test使用 若為true則getValue改吃mockMap的資料
	IsTest bool = false
	// Unit Test使用 若isTest為true則採用此資料
	MockMap map[string]string = make(map[string]string)
)

func GetCli() *agollo.Client {
	return cli
}

// 確定有呼叫InitApollo用才取得此結構體
type Listener struct{}

/* apollo setting by params */
func InitApollo(appid, host, namespace, cluster, secretkey string) Listener {
	agollo.SetLogger(zlog.GetSugarLog())
	c := &config.AppConfig{
		AppID:          appid,
		Cluster:        cluster,
		IP:             host,
		NamespaceName:  namespace,
		IsBackupConfig: true,
		Secret:         secretkey,
	}

	// 給予全域Namespace值
	Namespace = &namespace

	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	cli = client
	lis := Listener{}
	cli.AddChangeListener(lis)
	return lis
}

/* init with local yml file */
func InitWithConfig(ymlFile ...string) Listener {
	// get configMap data from yaml
	var configMap map[string]interface{}
	if len(ymlFile) > 0 {
		v := viper.New()
		v.SetConfigFile(ymlFile[0])
		v.SetConfigType("yaml")
		err := v.ReadInConfig()
		if err == nil {
			configMap = v.GetStringMap("apollo")
		} else {
			fmt.Println("apollo read yml file fail", err)
		}
	}
	// setting value
	appid := flag.String("appid", "", "")
	host := flag.String("host", "", "")
	namespace := flag.String("namespace", "", "")
	cluster := flag.String("cluster", "", "")
	secretkey := flag.String("secretkey", "", "")
	flag.Parse()

	// if nil , catch env or config
	appid = getConfig(appid, "appid", configMap)
	host = getConfig(host, "host", configMap)
	namespace = getConfig(namespace, "namespace", configMap)
	cluster = getConfig(cluster, "cluster", configMap)
	secretkey = getConfig(secretkey, "secretkey", configMap)

	fmt.Printf("初始化apollo設定 appid=%v,host=%v,namespace=%v,cluster=%v,secretkey=%v\n",
		*appid, *host, *namespace, *cluster, *secretkey)
	return InitApollo(*appid, *host, *namespace, *cluster, *secretkey)
}

/*
	參數取得設定檔
	優先級順序 flag > env > yml
*/
func getConfig(s *string, key string, ymlCfg map[string]interface{}) *string {
	// get from env
	if *s == "" && os.Getenv(key) != "" {
		*s = os.Getenv(key)
	}

	// get from config.yml
	if *s == "" && ymlCfg != nil && ymlCfg[key] != nil {
		*s = ymlCfg[key].(string)
	}

	return s
}

/*
	get config value from apollo
	must call InitApollo or InitWithConfig before this
*/
func GetValue(key string) (string, error) {
	if IsTest {
		return MockMap[key], nil
	}
	if cli == nil {
		return "", errors.New("no setting cli")
	}
	cache := cli.GetConfigCache(*Namespace)
	if cache == nil {
		return "", errors.New("get apollo cache fail")
	}
	value, err := cache.Get(key)
	if value != nil {
		str := utils.ToStr(value)
		return str, err
	}
	return cli.GetValue(key), err
}

/*
	get yaml setting value from apollo
*/
func GetYmlValue(value string) (map[string]interface{}, error) {
	// 從cli取得string
	str, err := GetValue(value)
	if err != nil {
		return nil, err
	}
	// yaml解析
	var m map[string]interface{}
	err = yaml.Unmarshal([]byte(str), &m)
	return m, err
}

/*
	bind yaml setting value from apollo
*/
func BindYmlValue(value string, val interface{}) error {
	// 從cli取得string
	str, err := GetValue(value)
	if err != nil {
		return err
	}
	// yaml解析
	err = yaml.Unmarshal([]byte(str), val)
	return err
}

/* 取得json對應的map值 */
func GetJsonValue(key string) (map[string]interface{}, error) {
	str, err := GetValue(key)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal([]byte(str), &m)
	return m, err
}

/* 從設定檔取得JSON值並綁定到obj */
func BindJsonValue(key string, val interface{}) error {
	str, err := GetValue(key)
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(str), val)
	return err
}
