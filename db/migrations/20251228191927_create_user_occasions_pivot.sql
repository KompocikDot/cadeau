-- migrate:up
CREATE TABLE IF NOT EXISTS user_occasions(
	occasion_id INTEGER NOT NULL REFERENCES occasions(id),
	user_id INTEGER NOT NULL REFERENCES users(id)
); 

-- migrate:down
DROP TABLE IF EXISTS user_occasions; 
