package credis

import (
	"context"
	"log"
	"time"

	redis "github.com/go-redis/redis/v8"
	redsync "github.com/go-redsync/redsync/v4"
	goredis "github.com/go-redsync/redsync/v4/redis/goredis/v8"
	apmgoredis "go.elastic.co/apm/module/apmgoredisv8"
)

var client *redis.Client

var clusterClient *redis.ClusterClient

var isCluster bool = false

func GetInstance() redis.Cmdable {
	if isCluster {
		return clusterClient
	} else {
		return client
	}
}

func NewRedis(addr string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     "",
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	rdb.AddHook(apmgoredis.NewHook())

	// 保存到全域變數
	client = rdb

	// 確認連線正常
	if _, err := rdb.Ping(context.TODO()).Result(); err != nil {
		log.Panic(err)
	}

	// 創建redsync
	pool := goredis.NewPool(rdb)
	rs = redsync.New(pool)
	isCluster = false

	return rdb
}

func SetRedis(r *redis.Client) {
	// 保存到全域變數
	client = r

	// 創建redsync
	pool := goredis.NewPool(r)
	rs = redsync.New(pool)
	isCluster = false
}

func NewRedisCluster(addrs []string) *redis.ClusterClient {
	rcdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Password:     "",
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})
	rcdb.AddHook(apmgoredis.NewHook())

	// 保存到全域變數
	clusterClient = rcdb

	if _, err := rcdb.Ping(context.TODO()).Result(); err != nil {
		log.Panic(err)
	}

	// 創建redsync
	pool := goredis.NewPool(rcdb)
	rs = redsync.New(pool)
	isCluster = true

	return rcdb
}

func SetRedisCluster(r *redis.ClusterClient) {
	// 保存到全域變數
	clusterClient = r

	// 創建redsync
	pool := goredis.NewPool(r)
	rs = redsync.New(pool)
	isCluster = true
}
