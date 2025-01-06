package dao

import (
	"fmt"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"

	"github.com/SongZihuan/chuangxinyi-communicate/model"
)

var (
	db *gorm.DB
)

const DRIVER_MYSQL = "mysql"

// Setup : Connect to mysql database
func Setup() errors.WTError {
	var err error

	switch viper.GetString("database.driver") {
	case DRIVER_MYSQL:
		host := viper.GetString("database.mysql.host")
		user := viper.GetString("database.mysql.user")
		password := viper.GetString("database.mysql.password")
		name := viper.GetString("database.mysql.name")
		charset := viper.GetString("database.mysql.charset")

		dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local", user, password, host, name, charset)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: gormlogger.Default.LogMode(gormlogger.Silent),
		})
		if err != nil {
			return errors.WarpQuick(err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			return errors.WarpQuick(err)
		}

		sqlDB.SetMaxIdleConns(viper.GetInt("database.mysql.pool.min"))
		sqlDB.SetMaxOpenConns(viper.GetInt("database.mysql.pool.max"))
		sqlDB.SetConnMaxLifetime(time.Minute)
	default:
		return errors.Errorf("we do not support this kind of storage system yet")
	}

	if err = db.AutoMigrate(model.Models...); nil != err {
		return errors.WarpQuick(err)
	}

	return nil
}

// Shutdown - close database connection
func Shutdown() errors.WTError {
	// 不需要处理关闭数据库
	return nil
}

// GetDb - get a database connection
func DB() *gorm.DB {
	return db
}

// 事务环绕
func Tx(db *gorm.DB, txFunc func(tx *gorm.DB) errors.WTError) (err errors.WTError) {
	tx := db.Begin()
	if tx.Error != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic after Rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = errors.WarpQuick(tx.Commit().Error)
		}
	}()

	return txFunc(tx)
}
