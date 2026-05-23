INSERT INTO `run`
(`inbox`, `started_at`, `finished_at`, `message_count`, `skipped_count`, `failed_count`, `moved_count`)
VALUES
(:inbox, :started_at, :finished_at, :message_count, :skipped_count, :failed_count, :moved_count)
;