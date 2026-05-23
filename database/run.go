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

type RunSummary struct {
	Inbox        string `db:"inbox"         json:"inbox"`
	RunCount     int    `db:"run_count"     json:"run_count"`
	FirstRunAt   string `db:"first_run_at"  json:"first_run_at"`
	LastRunAt    string `db:"last_run_at"   json:"last_run_at"`
	MessageCount int    `db:"message_count" json:"message_count"`
	SkippedCount int    `db:"skipped_count" json:"skipped_count"`
	FailedCount  int    `db:"failed_count"  json:"failed_count"`
	MovedCount   int    `db:"moved_count"   json:"moved_count"`
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

func ListRunSummaries(inbox string) ([]RunSummary, error) {
	var summaries []RunSummary
	if err := db.Select(&summaries, query("run/summary"), inbox, inbox); err != nil {
		return nil, err
	}
	return summaries, nil
}
