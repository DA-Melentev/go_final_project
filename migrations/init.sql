CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT,
    title TEXT,
    comment TEXT,
    repeat TEXT
);

CREATE INDEX IF NOT EXISTS scheduler_date_index ON scheduler(date);