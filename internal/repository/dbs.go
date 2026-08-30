package repository

import (
	"fmt"
	"log"
	"setaapi/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	DB_MEZZI   *sqlx.DB
	DB_CONTENT *sqlx.DB
	err        error
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
	DB_MEZZI, err = newDB(config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, "ertpl_mezzi")
	if err != nil {
		log.Fatal("Connessione DB fallita: ", err)
	}
}

func InitContent() {
	//DB connection
	DB_CONTENT, err = newDB(config.DB_HOST, config.DB_PORT, config.DB_USER, config.DB_PASS, "seta_api_content")
	if err != nil {
		log.Fatal("Connessione DB fallita: ", err)
	}
}
