package capollo

import (
	"fmt"
	"strings"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/rickylin614/common/cetcd"
	"github.com/rickylin614/common/cgorm"
	"github.com/rickylin614/common/ckafka"
	"github.com/rickylin614/common/credis"
	"github.com/rickylin614/common/utils"
	"github.com/rickylin614/common/zlog"
)

var ExecFuncs map[string]func(string) = make(map[string]func(string))

func init() {
	ExecFuncs["mysql"] = DbAutoSetting
	ExecFuncs["kafka"] = KafkaSet
	ExecFuncs["rediscluster"] = RedisClusterSet
	ExecFuncs["redis"] = RedisSet
	ExecFuncs["log"] = LogSet
	ExecFuncs["etcdRegister"] = EtcdRegister
	ExecFuncs["etcdClient"] = EtcdClient
}

/* apollo變更時 跟著變更 */
func (this Listener) OnChange(event *storage.ChangeEvent) {
	zlog.Info("apollo onchange execute")
	// 執行共用自動設定
	this.AutoSetting()
	// 執行被指定的重新初始化項目0
	for key := range event.Changes {
		if _, ok := ExecFuncs[key]; ok {
			zlog.Info("change key:", key)
			ExecFuncs[key](key)
		}
	}
}

/* 目前無作用 */
func (Listener) OnNewestChange(event *storage.FullChangeEvent) {}

/* 修改/添加onchange時執行的func */
func AddChangeFunc(key string, f func(string)) {
	ExecFuncs[key] = f
}

/* 修改/添加onchange時執行的funcs */
func (this Listener) AddChangeFuncs(fs map[string]func(string)) {
	for key, value := range fs {
		ExecFuncs[key] = value
	}
}

/* 根據support自動設定連線源 */
func (this Listener) AutoSetting() {
	// 從support取得欲初始化的項目
	settingStr, err := GetValue("support")
	var settingList []string
	if err == nil && settingStr != "" {
		settingList = strings.Split(settingStr, ",")
	} else {
		fmt.Println("auto setting err:", err)
		return
	}
	// 先做log
	if ExecFuncs["log"] != nil {
		ExecFuncs["log"]("log")
	}
	// 做log以外
	for _, v := range settingList {
		if ExecFuncs[v] != nil {
			ExecFuncs[v](v)
		}
	}
}

/* 靠apollo設定初始化DB設定 */
func DbAutoSetting(key string) {
	defer utils.ErrRecover()
	m, _ := GetYmlValue(key)
	confList := m["mysql"]
	if conf, ok := confList.([]interface{}); ok {
		for _, v := range conf {
			if confMap, ok := v.(map[string]interface{}); ok {
				if confMap["host"] == nil && confMap["schema"] == nil &&
					confMap["user"] == nil && confMap["pwd"] == nil {
					return
				}
				host := confMap["host"].(string)
				schema := confMap["schema"].(string)
				user := confMap["user"].(string)
				pwd := confMap["pwd"].(string)
				source := ""
				if confMap["source"] != nil {
					source = confMap["source"].(string)
				}
				cgorm.InitDB(host, schema, user, pwd, source)
			}
		}
	}
}

/* kafka連線設定 */
func KafkaSet(key string) {
	defer utils.ErrRecover()
	m, _ := GetYmlValue(key)
	if s, ok := m["brokers"].(string); ok {
		ckafka.Manage.SetBrokers(strings.Split(s, ","))
	}
	if s, ok := m["leader"].(string); ok {
		ckafka.Manage.SetLeaderAddr(s)
	}
}

/* 靠apollo設定初始化redisCluster設定 */
func RedisClusterSet(key string) {
	defer utils.ErrRecover()
	m, _ := GetYmlValue(key)
	if m["host"] == nil {
		return
	}
	var iList []interface{} = m["host"].([]interface{})
	hosts := make([]string, 0)
	for _, port := range iList {
		if v, ok := port.(string); ok {
			hosts = append(hosts, v)
		}
	}
	if len(hosts) != 0 {
		credis.NewRedisCluster(hosts)
	}
}

/* 靠apollo設定初始化redis設定 */
func RedisSet(key string) {
	defer utils.ErrRecover()
	m, _ := GetYmlValue(key)
	if m["host"] == nil {
		return
	}
	var host string = m["host"].(string)

	credis.NewRedis(host)
}

/* 靠apollo設定初始化LogSet設定 */
func LogSet(key string) {
	defer utils.ErrRecover()
	m, _ := GetYmlValue(key)
	if m["infopath"] == nil || m["errorpath"] == nil {
		zlog.ConsoleInit()
		return
	}
	zlog.InitLog(m["infopath"].(string), m["errorpath"].(string), true)
}

/* apollo設定初始化etcd Register設定 */
func EtcdRegister(key string) {
	defer utils.ErrRecover()
	m, _ := GetYmlValue(key)
	host, err := GetValue("host")
	if err != nil {
		fmt.Println("etcd register get host err:", err)
	}
	if host == "" {
		host = "127.0.0.1:8080"
	}
	if m["prefix"] == nil || m["endpoints"] == nil {
		return
	}
	if s, ok := m["endpoints"].(string); ok {
		endpointList := strings.Split(s, ",")
		srv, err := cetcd.NewService(m["prefix"].(string), host, endpointList)
		if err != nil {
			fmt.Println("etcd register err:", err)
		}
		go srv.Start()
	}
}

/* apollo設定初始化etcd 服務發現用設定 */
func EtcdClient(key string) {
	defer utils.ErrRecover()
	m, _ := GetYmlValue(key)
	if s, ok := m["endpoints"].(string); ok {
		endpointList := strings.Split(s, ",")
		cetcd.NewClientDis(endpointList)
	}
}
