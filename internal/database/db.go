package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"string_backend_0001/internal/conf"
	logger2 "string_backend_0001/internal/logger"
	"strings"
	"time"
)

var db *gorm.DB

func Init() error {
	var err error
	database := conf.Conf.Database

	newLogger := logger.New(logger2.GetLogger(), logger.Config{
		SlowThreshold:             time.Second,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	})

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.Conf.Database.TablePrefix,
		},
		Logger: newLogger,
	}
	switch database.Type {
	case "sqlite3":
		if !(strings.HasSuffix(database.DBFile, ".db") && len(database.DBFile) > 3) {
			log.Fatalf("db name error.")
		}
		db, err = gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",
			database.DBFile)), gormConfig)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%s",
			database.User, database.Password, database.Host, database.Port, database.Name, database.SSLMode)
		if database.DSN != "" {
			dsn = database.DSN
		}
		db, err = gorm.Open(mysql.Open(dsn), gormConfig)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
			database.Host, database.User, database.Password, database.Name, database.Port, database.SSLMode)
		if database.DSN != "" {
			dsn = database.DSN
		}
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	default:
		err = fmt.Errorf("not supported database type: %s", database.Type)
	}

	return err
}

func Close() error {
	if db == nil {
		return nil
	}

	s, err := db.DB()
	if err != nil {
		return err
	}

	return s.Close()
}

func GetDB() *gorm.DB {
	return db
}
