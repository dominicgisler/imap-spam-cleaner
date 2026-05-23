package database

import (
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/pressly/goose/v3"
)

const (
	migVersion = 1
	migDir     = "queries/migrations"
)

type gooseLogger struct{}

func (*gooseLogger) Fatalf(format string, v ...interface{}) {
	logx.Panicf(format, v...)
}

func (*gooseLogger) Printf(format string, v ...interface{}) {
	logx.Infof(format, v...)
}

func migrateDB() error {

	goose.SetLogger(&gooseLogger{})
	goose.SetBaseFS(fs)

	if err := goose.SetDialect(driverName); err != nil {
		return err
	}

	ver, err := goose.GetDBVersion(db.DB)
	if err != nil {
		return err
	}

	if ver > migVersion {
		return goose.DownTo(db.DB, migDir, migVersion)
	}
	return goose.UpTo(db.DB, migDir, migVersion)
}
