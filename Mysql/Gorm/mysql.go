package Gorm

import (
	"errors"
	"fmt"
	"github.com/chaos-star/marvel/Env"
	"github.com/chaos-star/marvel/Log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"log"
	"os"
	"time"
)

type Engine struct {
	db map[string]*gorm.DB
}

func (e *Engine) DB(name string) *gorm.DB {
	if db, ok := e.db[name]; ok {
		return db
	}
	return nil
}
func (e *Engine) Instance(name string) *gorm.DB {
	if db, ok := e.db[name]; ok {
		return db
	}
	return nil
}
func (e *Engine) DbMap() map[string]*gorm.DB {
	return e.db
}

type mysqlOption struct {
	dbLog       bool
	charset     string
	maxIdleConn int
	maxOpenConn int
	maxLifetime time.Duration
	maxIdleTime time.Duration
}

type mysqlConfig struct {
	dsn    string
	option mysqlOption
	slave  []*mysqlConfig
}

func Initialize(env string, mysqlConfigs interface{}, mLog Log.ILogger) (*Engine, error) {
	var (
		dbs = make(map[string]*gorm.DB)
		mcs []map[string]interface{}
	)
	if v, ok := mysqlConfigs.(map[string]interface{}); ok {
		mcs = append(mcs, v)
	}

	if v, ok := mysqlConfigs.([]interface{}); ok {
		if len(v) > 0 {
			for _, iv := range v {
				if ivm, is := iv.(map[string]interface{}); is {
					mcs = append(mcs, ivm)
				}
			}
		}
	}

	if len(mcs) > 0 {
		for _, mc := range mcs {
			db, err := newDB(env, mc, mLog)
			if err != nil {
				return nil, err
			}
			dbs[mc["database"].(string)] = db
			if alias, has := mc["alias"]; has && alias != "" {
				dbs[alias.(string)] = db
			}
		}
	}

	return &Engine{dbs}, nil
}

func newDB(env string, conf map[string]interface{}, iLog Log.ILogger) (*gorm.DB, error) {
	mc, err := parseDBConfig(conf)
	if err != nil {
		return nil, err
	}
	var newLogger logger.Interface
	if mc.option.dbLog && env != Env.DeployEnvProd && env != Env.DeployEnvDebug {
		newLogger = logger.New(
			iLog,
			logger.Config{
				SlowThreshold:             time.Millisecond * 20, // Slow SQL threshold
				LogLevel:                  logger.Info,           // Log level
				IgnoreRecordNotFoundError: false,                 // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                  // Disable color
			},
		)
	} else {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Millisecond * 200, // Slow SQL threshold
				LogLevel:                  logger.Info,            // Log level
				IgnoreRecordNotFoundError: false,                  // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                   // Disable color
			},
		)
	}

	return initDB(mc, newLogger)
}

func initDB(conf *mysqlConfig, iLog logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conf.dsn, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,      // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,     // disable datetime precision support, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,     // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,     // use change when rename column, rename rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,    // smart configure based on used version
	}), &gorm.Config{
		Logger: iLog,
	})

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if conf.option.maxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(conf.option.maxOpenConn)
	}
	if conf.option.maxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(conf.option.maxIdleConn)
	}
	if conf.option.maxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(conf.option.maxIdleTime)
	}
	if conf.option.maxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(conf.option.maxLifetime)
	}

	if len(conf.slave) > 0 {
		var Replicas []gorm.Dialector
		for _, v := range conf.slave {
			Replicas = append(Replicas, mysql.Open(v.dsn))
		}
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(conf.dsn)},
			Replicas: Replicas,
			Policy:   dbresolver.RandomPolicy{},
		}))
		if err != nil {
			return nil, err
		}
	}
	return db, err
}

func parseDBConfig(conf map[string]interface{}) (*mysqlConfig, error) {
	if conf["host"] == nil || conf["host"] == "" {
		return nil, errors.New("host is not found")
	}
	if conf["username"] == nil || conf["username"] == "" {
		return nil, errors.New("username is not found")
	}
	if conf["password"] == nil || conf["password"] == "" {
		return nil, errors.New("password is not found")
	}
	if conf["database"] == nil || conf["database"] == "" {
		return nil, errors.New("database is not found")
	}
	var (
		option mysqlOption
		mSlave []map[string]interface{}
		slave  []*mysqlConfig
	)
	option.maxIdleConn = 10
	option.maxOpenConn = 100
	option.maxLifetime = time.Hour
	option.maxIdleTime = time.Minute * 15
	option.charset = "utf8mb4"

	if maxOpen, ok := conf["max_open"]; ok {
		if val, is := maxOpen.(int); is && val > 0 {
			option.maxOpenConn = val
		}
	}

	if maxIdle, ok := conf["max_idle"]; ok {
		if val, is := maxIdle.(int); is && val > 0 {
			option.maxIdleConn = val
		}
	}

	if maxLifeTime, ok := conf["max_life_time"]; ok {
		if val, is := maxLifeTime.(time.Duration); is && val > 0 {
			option.maxLifetime = val * time.Minute
		}
	}

	if maxIdleTime, ok := conf["max_idle_time"]; ok {
		if val, is := maxIdleTime.(time.Duration); is && val > 0 {
			option.maxIdleTime = val * time.Minute
		}
	}

	if charset, ok := conf["charset"]; ok {
		if val, is := charset.(string); is && val != "" {
			option.charset = val
		}
	}
	if dbLog, ok := conf["db_log"]; ok {
		if val, is := dbLog.(bool); is {
			option.dbLog = val
		}
	}
	conf["max_open"] = option.maxOpenConn
	conf["max_idle"] = option.maxIdleConn
	conf["max_life_time"] = option.maxLifetime
	conf["max_idle_time"] = option.maxIdleTime
	conf["charset"] = option.charset
	conf["db_log"] = option.dbLog

	if sc, ok := conf["slave"].(map[string]interface{}); ok {
		mSlave = append(mSlave, sc)
	}

	if sc, ok := conf["slave"].([]interface{}); ok {
		if len(sc) > 0 {
			for _, isc := range sc {
				if iscm, is := isc.(map[string]interface{}); is {
					mSlave = append(mSlave, iscm)
				}
			}
		}
	}

	if len(mSlave) > 0 {
		for _, msc := range mSlave {
			s, err := parseDBConfig(msc)
			if err != nil {
				return nil, err
			}
			slave = append(slave, s)
		}
	}

	return &mysqlConfig{
		dsn:    parseDSN(conf),
		option: option,
		slave:  slave,
	}, nil
}

func parseDSN(conf map[string]interface{}) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", conf["username"], conf["password"], conf["host"], conf["database"], conf["charset"])
	return dsn
}
