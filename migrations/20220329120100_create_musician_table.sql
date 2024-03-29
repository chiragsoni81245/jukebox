-- +goose Up
CREATE TABLE musicians (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    musician_type TEXT
);

-- +goose Down
DROP TABLE musicians;

