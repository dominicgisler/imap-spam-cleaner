-- +goose Up
-- +goose StatementBegin
CREATE TABLE `run`
(
    `id`            INTEGER PRIMARY KEY AUTOINCREMENT,
    `inbox`         TEXT      NOT NULL,
    `started_at`    TIMESTAMP NOT NULL,
    `finished_at`   TIMESTAMP,
    `message_count` INTEGER   NOT NULL DEFAULT 0,
    `skipped_count` INTEGER   NOT NULL DEFAULT 0,
    `failed_count`  INTEGER   NOT NULL DEFAULT 0,
    `moved_count`   INTEGER   NOT NULL DEFAULT 0
);

CREATE INDEX `idx_run_inbox` ON `run` (`inbox`);
CREATE INDEX `idx_run_inbox_started_at` ON `run` (`inbox`, `started_at`);
CREATE INDEX `idx_run_started_at` ON `run` (`started_at`);
CREATE INDEX `idx_run_finished_at` ON `run` (`finished_at`);
CREATE INDEX `idx_run_counts` ON `run` (`message_count`, `skipped_count`, `failed_count`, `moved_count`);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE `run`;
DROP INDEX `idx_run_counts`;
DROP INDEX `idx_run_finished_at`;
DROP INDEX `idx_run_started_at`;
DROP INDEX `idx_run_inbox_started_at`;
-- +goose StatementEnd
