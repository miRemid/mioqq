package http

import (
	"time"

	"github.com/jinzhu/gorm"
	// mssql 微软数据库
	_ "github.com/jinzhu/gorm/dialects/mssql"
	// mysql 数据库
	_ "github.com/jinzhu/gorm/dialects/mysql"
	// postgres data
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// sqlite data
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	database *gorm.DB
)

type noticeMessage struct {
	gorm.Model
}

// SetDatabase 设置数据库连接
func (server *Server) SetDatabase(driver, uri string) error {
	server.Driver = driver
	server.ConnectURI = uri
	return server.connect()
}

func (server *Server) connect() error {
	db, err := gorm.Open(server.Driver, server.ConnectURI)
	if err != nil {
		return err
	}
	database = db
	database.SingularTable(true)
	database.DB().SetConnMaxLifetime(time.Second * 30)
	database.DB().SetMaxIdleConns(20)
	database.DB().SetMaxOpenConns(100)

	database.AutoMigrate(&noticeMessage{})

	return nil
}
