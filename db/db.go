package db

import (
	"fmt"
	"hand_held/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var (
	DBConn *gorm.DB
)

func InitDatabase() error {
	var err error
	connectionString := fmt.Sprintf("sqlserver://%s:%s@%s:1433?database=%s", config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_HOST"), config.Config("DB_NAME"))
	fmt.Println(connectionString)
	// DBConn, err = gorm.Open("mssql", "sqlserver://mcs:123@41.38.87.59:1433?database=stock_main")
	DBConn, err = gorm.Open("mssql", connectionString)
	if err != nil {
		fmt.Println("Failed to connect to external database")
		panic(err)
	}
	DBConn.LogMode(true)
	fmt.Println("Connection Opened to External Database")
	return nil

}
