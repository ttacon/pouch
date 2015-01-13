package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ttacon/pouch/impl"
	"github.com/ttacon/pretty"
)

var (
	// for DB mode
	host     = flag.String("h", "", "db host to connect to")
	user     = flag.String("u", "", "user to connect to db as")
	password = flag.String("p", "", "password to authenticate user with")
	database = flag.String("db", "", "database to read tables from")
)

func main() {
	flag.Parse()

	dbConn, err := getDBConn(*host, *user, *password, *database)
	if err != nil {
		fmt.Println("failed to connect to the db, err: ", err)
		return
	}

	p := impl.SQLPouch(dbConn)
	var f Food = Food{ID: 1}
	err = p.Find(&f)
	if err != nil {
		fmt.Println("failed to find Food: ", err)
		return
	}

	pretty.Println(f)
}

func getDBConn(host, username, password, database string) (*sql.DB, error) {
	return sql.Open("mysql",
		fmt.Sprintf("%s:%s@%s/%s?parseTime=true", username, password, host, database))
}
