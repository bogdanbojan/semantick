-- Version: 1.1
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID,
	name          TEXT,
	email         TEXT UNIQUE,
	roles         TEXT[],
	password_hash TEXT,
	date_created  TIMESTAMP,
	date_updated  TIMESTAMP,

	PRIMARY KEY (user_id)
);


-- Version: 1.2
-- Description: Create table books
CREATE TABLE books (
	book_id      UUID,
	name         TEXT,
	highlights   INT,
	user_id      UUID,
	date_created TIMESTAMP,
	date_updated TIMESTAMP,

	PRIMARY KEY (book_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);


-- Version: 1.3
-- Description: Create table highlights
CREATE TABLE highlights (
	highlight_id         UUID,
	user_id              UUID,
	book_id              UUID,
	definition           TEXT,
	definition_provider  TEXT,
	date_created         TIMESTAMP,
	date_updated         TIMESTAMP,

	PRIMARY KEY (book_id),
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
	FOREIGN KEY (highlight_id) REFERENCES books(book_id) ON DELETE CASCADE
);
