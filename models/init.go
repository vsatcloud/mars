package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

type Database struct {
	Host         string `json:"host" required:"true" env:"PGDB_ADDR"`
	Port         string `json:"port" required:"true" env:"PGDB_PORT"`
	User         string `json:"user" required:"true" env:"PGDB_USER"`
	Dbname       string `json:"dbname" required:"true" env:"PGDB_DB_NAME"`
	Password     string `json:"password" required:"true" env:"PGDB_PASSWORD"`
	Sslmode      string `json:"sslmode" required:"true"`
	MaxIdleConns int    `json:"max_idle_conns" required:"true"`
	MaxOpenConns int    `json:"max_open_conns" required:"true"`
}

func envDefault(key, defaultKey string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultKey
}

func Init(conf Database) error {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		envDefault("PGDB_HOST", conf.Host),
		envDefault("PGDB_PORT", conf.Port),
		envDefault("PGDB_USER", conf.User),
		envDefault("PGDB_DB_NAME", conf.Dbname),
		envDefault("PGDB_PASSWORD", conf.Password),
		conf.Sslmode)
	db, err = gorm.Open("postgres", connStr)
	if err != nil {
		return err
	}

	var modelsList = []interface{}{OperationRecord{}}

	if err = db.AutoMigrate(modelsList...).Error; nil != err {
		return err
	}

	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)

	return nil
}

// LogMode set log mode, `true` for detailed logs, `false` for no log, default, will only print error logs
func LogMode(enable bool) {
	db.LogMode(enable)
}

func CloseDB() error {
	return db.Close()
}

func GetDB() *gorm.DB {
	return db
}

//Offset 获取数据库查询的offset
func Offset(page, limit int) int {
	if page <= 0 {
		return -1 //cancel offset
	}

	return (page - 1) * limit
}

func Limit(limit int) int {
	if limit <= 0 {
		return -1
	}

	return limit
}

//自动构造表结构
func AutoMigrateModel(values ...interface{}) error {
	if err := db.AutoMigrate(values).Error; nil != err {
		return err
	}

	return nil
}
