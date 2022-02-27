create table log
(
    date_time DateTime,
    user_id   Int64
) ENGINE = MergeTree()
      partition by toYYYYMMDD(date_time)
      order by (date_time);