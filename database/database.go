package database

import (
	"embed"
	"fmt"

	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

const driverName = "sqlite3"

//go:embed queries/*
var fs embed.FS

func Init(path string) (err error) {

	if path == "" {
		return fmt.Errorf("database path is required")
	}

	db, err = sqlx.Open(driverName, path)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping database: %w", err)
	}

	if err = migrateDB(); err != nil {
		_ = db.Close()
		return fmt.Errorf("migrate database: %w", err)
	}

	return nil
}

func query(file string) string {
	bs, err := fs.ReadFile("queries/" + file + ".sql")
	if err != nil {
		logx.Panic(err.Error())
	}
	return string(bs)
}
