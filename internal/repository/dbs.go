package repository

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB_MEZZI *sqlx.DB
	err      error
)

func newDB(host string, port int, user, password, dbname string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", user, password, host, port, dbname)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func InitMezzi() {
	//DB connection
	//NOTE FOR ME CHANGE PASSWORD WHEN PRODUCTION!!!
	DB_MEZZI, err = newDB("serverissimo.com", 3306, "setaapi", "71070c80767dbdc02cbbeca9ed9841b4", "ertpl_mezzi")
	if err != nil {
		log.Fatal("Connessione DB fallita: ", err)
	}
	//defer DB_MEZZI.Close()
}
