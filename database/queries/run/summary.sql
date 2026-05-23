SELECT
    `inbox`,
    COUNT(*) AS `run_count`,
    MIN(`started_at`) AS `first_run_at`,
    MAX(`started_at`) AS `last_run_at`,
    SUM(`message_count`) AS `message_count`,
    SUM(`skipped_count`) AS `skipped_count`,
    SUM(`failed_count`) AS `failed_count`,
    SUM(`moved_count`) AS `moved_count`
FROM `run`
WHERE (? = '' OR `inbox` = ?)
GROUP BY `inbox`
ORDER BY `inbox`;
