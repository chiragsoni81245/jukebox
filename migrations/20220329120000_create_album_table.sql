-- +goose Up
CREATE TABLE albums (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    release_date TIMESTAMP NOT NULL,
    genre TEXT,
    price REAL NOT NULL,
    description TEXT
);

-- +goose Down
DROP TABLE albums;

