func initDB() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "./bot.db")
	if err != nil {
		return nil, err
	}

	// Create Users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		userID INTEGER PRIMARY KEY AUTOINCREMENT,
		telegramID INTEGER UNIQUE,
		username VARCHAR(255),
		points INTEGER DEFAULT 0,
		alias VARCHAR(255) DEFAULT ''
	)`)
	if err != nil {
		return nil, err
	}

	// Create Tasks table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (
		taskID INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		description TEXT DEFAULT '',
		points INTEGER DEFAULT 0,
		periodicity INTEGER DEFAULT 7,
		status TEXT DEFAULT 'waiting',
		doneLast DATETIME DEFAULT "0000-00-00 00:00:00"
	)`)
	if err != nil {
		return nil, err
	}

	// Create Actions table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS actions (
		actionID INTEGER PRIMARY KEY AUTOINCREMENT,
		taskID INTEGER NOT NULL,
		userID INTEGER,
		timestamp DATETIME,
		FOREIGN KEY (taskID) REFERENCES tasks(taskID),
		FOREIGN KEY (userID) REFERENCES users(userID)
	)`)
	if err != nil {
		return nil, err
	}

	// Create messages table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS allMessages
		(id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		username TEXT,
		message_text TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		log.Panic(err)
	}

	log.Println("database created")
	return db, nil
}