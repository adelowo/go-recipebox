CREATE TABLE IF NOT EXISTS users (

	id integer primary key autoincrement,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password VARCHAR(300) NOT NULL
);

CREATE TABLE IF NOT EXISTS recipes (

	id integer primary key autoincrement,
	user_id int  NOT NULL,
	title VARCHAR(300) NOT NULL,
	ingredients TEXT NOT NULL,
	description TEXT NOT NULL
);
