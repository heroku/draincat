package main

import (
	"database/sql"
	"log"
)

var db *sql.DB
var insertStmt *sql.Stmt

func Insert(privalVersion, time, hostname, name, procid, msgid, data []byte) error {
	_, err := insertStmt.Exec(
		string(privalVersion),
		string(time),
		string(hostname),
		string(name),
		string(procid),
		string(msgid),
		string(data))
	return err
}

func init() {
	// Connect to postgresql
	var err error
	db, err = sql.Open("postgres", config.DatabaseUrl)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to ping DB: %v", err)
	}

	insertStmt, err = db.Prepare(
		"INSERT into logs(privalversion, time, hostname, name, procid, msgid, data) VALUES ($1, $2, $3, $4, $5, $6, $7);")
	if err != nil {
		log.Fatalf("Unable to create prepared stmt: %v", err)
	}
}
