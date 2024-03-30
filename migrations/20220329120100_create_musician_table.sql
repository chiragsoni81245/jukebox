-- +goose Up
CREATE TABLE musicians (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT
);

-- +goose Down
DROP TABLE musicians;

