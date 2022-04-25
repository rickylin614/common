package cgorm

import (
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rickylin614/common/zlog"

	apmmysql "go.elastic.co/apm/module/apmgormv2/driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

// for single database
var db *gorm.DB

// for multi database
var dbs map[string]*gorm.DB = make(map[string]*gorm.DB)

/*
	host: host+port
	schema: schema名稱
	user: 使用者帳號
	password: 密碼
	dbSourceName: 多連線源使用 給予該連線名稱 若不使用則給空字串
*/
func InitDB(host, schema, user, password, dbSourceName string) (err error) {
	dsn := user + ":" + password + "@tcp(" + host + ")/" + schema + "?charset=utf8mb4&loc=Local&parseTime=true&timeout=30s"

	// log初始化設定
	logger := zapgorm2.Logger{
		ZapLogger:                 zlog.GetLog(),
		LogLevel:                  gormlogger.Info,
		SlowThreshold:             100 * time.Millisecond,
		SkipCallerLookup:          false,
		IgnoreRecordNotFoundError: false,
	}

	// 連線
	gormdb, err := gorm.Open(apmmysql.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		return err
	}

	// 設置連接池數據
	sqlDB, err := gormdb.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Ping()
	if err != nil {
		return err
	}
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(100)
	// 閒置連線的最大存在時間
	sqlDB.SetConnMaxIdleTime(time.Second * 25)
	// 連線的最大生存時間 確保連線可以被驅動安全關閉 官方建議小於五分鐘
	sqlDB.SetConnMaxLifetime(time.Second * 25)

	err = sqlDB.Ping()
	if err != nil {
		zlog.Error(err)
		return
	}

	if dbSourceName == "" {
		db = gormdb
	} else {
		dbs[dbSourceName] = gormdb
	}
	return err
}

/*
	避免任何程序對初始化後的DB做設定修改，連線後"只能"對db設定做"讀取"
	!important do not change db.config
*/
func GetDB(sourceName ...string) *gorm.DB {
	if len(sourceName) == 0 {
		return db
	}
	if dbs[sourceName[0]] == nil {
		zlog.Panic("can't find the source config :", sourceName[0])
	}
	return dbs[sourceName[0]]
}

/* 多連線源 給連線源名稱 */
func Begin(sourceName ...string) *gorm.DB {
	if len(sourceName) == 0 {
		return db.Begin()
	}
	if dbs[sourceName[0]] == nil {
		zlog.Panic("can't find the source config :", sourceName[0])
	}
	return dbs[sourceName[0]].Begin()
}

func GetMock() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	return gormDB, mock
}

func NewMockDb(sourceName ...string) sqlmock.Sqlmock {
	mockDB, mock, _ := sqlmock.New()
	gormMockDB, _ := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      mockDB,
	}), &gorm.Config{})
	if len(sourceName) == 0 || sourceName[0] == "" {
		db = gormMockDB
	} else {
		dbs[sourceName[0]] = gormMockDB
	}
	return mock
}
