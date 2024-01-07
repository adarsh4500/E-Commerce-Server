package connections

import (
	"Ecom/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var err error

func InitializeDB() error {

	conn_string := "user=" + config.DB_UserName + " password=" + config.DB_Password + " dbname=" + config.DB_Name + " sslmode=disable"

	if DB, err = sql.Open("postgres", conn_string); err != nil {
		fmt.Println(err)
		return err
	}

	if err = DB.Ping(); err != nil {
		fmt.Println(err)
		return err
	}

	log.Print("Connected to database")
	return nil
}