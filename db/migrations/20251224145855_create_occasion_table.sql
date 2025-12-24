-- migrate:up
CREATE TABLE IF NOT EXISTS occasions(
id INTEGER PRIMARY KEY AUTOINCREMENT,
name VARCHAR NOT NULL,
gift_receiver INTEGER NOT NULL REFERENCES users(id)
);

-- migrate:down
DROP TABLE IF EXISTS occasions;
