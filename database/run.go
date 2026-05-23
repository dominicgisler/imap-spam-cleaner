package database

import (
	"time"
)

type Run struct {
	ID           int       `db:"id"`
	Inbox        string    `db:"inbox"`
	StartedAt    time.Time `db:"started_at"`
	FinishedAt   time.Time `db:"finished_at"`
	MessageCount int       `db:"message_count"`
	SkippedCount int       `db:"skipped_count"`
	FailedCount  int       `db:"failed_count"`
	MovedCount   int       `db:"moved_count"`
}

func AddRun(run *Run) error {
	res, err := db.NamedExec(query("run/create"), run)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	run.ID = int(id)
	return nil
}
