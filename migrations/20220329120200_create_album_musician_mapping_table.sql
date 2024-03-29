-- +goose Up
CREATE TABLE album_musician (
    album_id INTEGER,
    musician_id INTEGER,
    PRIMARY KEY (album_id, musician_id),
    FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
    FOREIGN KEY (musician_id) REFERENCES musicians(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE album_musician;

