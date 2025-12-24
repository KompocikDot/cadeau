CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE users(
  id INTEGER PRIMARY KEY,
  username VARCHAR NOT NULL UNIQUE,
  password VARCHAR NOT NULL
);
CREATE TABLE occasions(
id INTEGER PRIMARY KEY AUTOINCREMENT,
name VARCHAR NOT NULL,
gift_receiver INTEGER REFERENCES users(id)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20251224124538'),
  ('20251224145855');
