package dbops

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConn *sql.DB
	err    error
)

func init() {
	fmt.Println("Entering conn.go init function...")
	dbConn, err = sql.Open("mysql", "pace:123@/sample")
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("dbConn open +%v\n", dbConn)
}
