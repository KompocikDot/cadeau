CREATE TABLE IF NOT EXISTS "schema_migrations" (version varchar(128) primary key);
CREATE TABLE users(
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username VARCHAR NOT NULL UNIQUE,
  password VARCHAR NOT NULL
);
CREATE TABLE occasions(
id INTEGER PRIMARY KEY AUTOINCREMENT,
name VARCHAR NOT NULL,
gift_receiver INTEGER NOT NULL REFERENCES users(id)
);
CREATE TABLE gifts(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR NOT NULL,
	url VARCHAR NOT NULL,
	occasion INTEGER NOT NULL REFERENCES occasions(id)
);
CREATE TABLE user_occasions(
	occasion_id INTEGER NOT NULL REFERENCES occasions(id),
	user_id INTEGER NOT NULL REFERENCES users(id)
);
-- Dbmate schema migrations
INSERT INTO "schema_migrations" (version) VALUES
  ('20251224124538'),
  ('20251224145855'),
  ('20251224211835'),
  ('20251228191927');
