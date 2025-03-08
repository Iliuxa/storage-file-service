CREATE TABLE IF NOT EXISTS storage_file
(
    id          INTEGER PRIMARY KEY,
    file_name   TEXT    NOT NULL,
    insert_date TEXT    NOT NULL,
    update_date TEXT,
    delete_date TEXT,
    file_path   TEXT    NOT NULL UNIQUE,
    file_hash   BLOB    NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_file_id ON storage_file (id);