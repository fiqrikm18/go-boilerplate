package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fiqrikm18/go-boilerplate/internal/model/dao"
	"github.com/fiqrikm18/go-boilerplate/pkg/lib"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbConnection struct {
	DB *gorm.DB
}

func NewDbConnection() (*DbConnection, error) {
	var dialector gorm.Dialector
	appMode := os.Getenv("APP_MODE")
	appConf, err := lib.LoadConfigFile()
	if err != nil {
		return nil, err
	}

	dbDriver := appConf.DBConf.Driver
	dbName := appConf.DBConf.Name
	dbUsername := appConf.DBConf.Username
	dbPassword := appConf.DBConf.Password
	dbHost := appConf.DBConf.Host
	dbPort := appConf.DBConf.Port

	if appMode != "" && appMode == "test" {
		dbDriver = appConf.DBConfTest.Driver
		dbName = appConf.DBConfTest.Name
		dbUsername = appConf.DBConfTest.Username
		dbPassword = appConf.DBConfTest.Password
		dbHost = appConf.DBConfTest.Host
		dbPort = appConf.DBConfTest.Port
	}

	switch dbDriver {
	case "psql":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUsername, dbPassword, dbName, dbPort)
		dialector = postgres.Open(dsn)
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbPort, dbName)
		dialector = mysql.Open(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", dbUsername, dbPassword, dbHost, dbPort, dbName)
		dialector = sqlserver.Open(dsn)
	}

	conn, err := gorm.Open(dialector, &gorm.Config{
		Logger:                                   dbLogger(),
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, err
	}

	runAutoMigrate(conn)

	return &DbConnection{
		DB: conn,
	}, nil
}

func runAutoMigrate(conn *gorm.DB) {
	conn.AutoMigrate(
		dao.User{},
	)
}

func dbLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
