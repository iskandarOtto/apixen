CREATE TABLE IF NOT EXISTS organizations (
	id SERIAL PRIMARY KEY,
	fullname VARCHAR(256) NOT NULL UNIQUE,
	shortname VARCHAR(50) NOT NULL UNIQUE,
	created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS comments (
	id BIGSERIAL PRIMARY KEY,
	comment TEXT NOT NULL,
	orgid INT REFERENCES organizations(id),
	created TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS deletedcomments (
	id BIGSERIAL PRIMARY KEY,
	comment TEXT NOT NULL,
	orgid INT REFERENCES organizations(id),
	created TIMESTAMP NOT NULL,
	deletedate TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS members (
	id SERIAL PRIMARY KEY,
	email VARCHAR(256) NOT NULL UNIQUE,
	username VARCHAR(100) NOT NULL UNIQUE,
	passhash VARCHAR(256) NOT NULL,
	avatarurl VARCHAR(256),
	followerno INT DEFAULT 0,
	followingno INT DEFAULT 0,
	orgid INT REFERENCES organizations(id)
);