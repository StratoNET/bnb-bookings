package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB holds whichever database is employed (e.g. MariaDB, Postgres etc) connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDBConn = 10
const maxIdleDBConn = 10
const maxDBLifetime = 3 * time.Minute

// ConnectSQL creates pool for currently used database (MariaDB/MySQL)
func ConnectSQL(dsn string) (*DB, error) {
	dp, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	dp.SetMaxOpenConns(maxOpenDBConn)
	dp.SetMaxIdleConns(maxIdleDBConn)
	dp.SetConnMaxLifetime(maxDBLifetime)

	dbConn.SQL = dp

	err = testDB(dp)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// NewDatabase creates a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

// testDB local function which attempts to ping database creating test connection
func testDB(dp *sql.DB) error {
	err := dp.Ping()
	if err != nil {
		return err
	}

	return nil
}
