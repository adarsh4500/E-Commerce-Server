package connections

import (
	"Ecom/config"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var DB *sql.DB
var err error

func InitializeDB() error {

	conn_string := "host=" + config.DB_Host + " port=" + config.DB_Port + " user=" + config.DB_UserName + " password=" + config.DB_Password + " dbname=" + config.DB_Name + " sslmode=" + config.DB_SSLMode

	if DB, err = sql.Open("postgres", conn_string); err != nil {
		fmt.Println(err)
		return err
	}

	if err = DB.Ping(); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
