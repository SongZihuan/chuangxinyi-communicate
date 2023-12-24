package dao

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/spf13/viper"
	errors "github.com/wuntsong-org/wterrors"
	"time"

	"gitee.com/wuntsong/chuangxinyi-communicate/model"
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
		db, err = gorm.Open("mysql", dsn)
		if err != nil {
			return errors.WarpQuick(err)
		}

		db.DB().SetMaxIdleConns(viper.GetInt("database.mysql.pool.min"))
		db.DB().SetMaxOpenConns(viper.GetInt("database.mysql.pool.max"))
		db.DB().SetConnMaxLifetime(time.Minute)
		if gin.Mode() != gin.ReleaseMode {
			db.LogMode(true)
		}
	default:
		return errors.Errorf("we do not support this kind of storage system yet")
	}

	db.SingularTable(true) //禁用表名复数
	if err = db.AutoMigrate(model.Models...).Error; nil != err {
		return errors.WarpQuick(err)
	}

	return nil
}

// Shutdown - close database connection
func Shutdown() errors.WTError {
	return errors.WarpQuick(db.Close())
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
