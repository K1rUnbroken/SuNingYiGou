package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db *sql.DB

// InitDB 初始化数据库
func InitDB() {
	dsn := "root:123456789@tcp(localhost:3306)/suning"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
