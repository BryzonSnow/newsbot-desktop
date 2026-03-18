package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB(filepath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	queryNoticias := `CREATE TABLE IF NOT EXISTS sent_articles (url TEXT PRIMARY KEY, sent_at DATETIME DEFAULT CURRENT_TIMESTAMP);`
	database.Exec(queryNoticias)

	queryConfig := `
	CREATE TABLE IF NOT EXISTS user_config (
		id INTEGER PRIMARY KEY CHECK (id = 1),
		phone TEXT,
		news_api_key TEXT,
		wa_api_key TEXT,
		global_topics TEXT,
		local_topics TEXT,
		interval_minutes INTEGER DEFAULT 30
	);`
	database.Exec(queryConfig)

	return database, nil
}

func GuardarConfig(db *sql.DB, phone, newsKey, waKey, global, local string, intervalo int) error {
	query := `
	INSERT INTO user_config (id, phone, news_api_key, wa_api_key, global_topics, local_topics, interval_minutes) 
	VALUES (1, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(id) DO UPDATE SET 
		phone=excluded.phone, 
		news_api_key=excluded.news_api_key, 
		wa_api_key=excluded.wa_api_key, 
		global_topics=excluded.global_topics, 
		local_topics=excluded.local_topics,
		interval_minutes=excluded.interval_minutes;`
	
	_, err := db.Exec(query, phone, newsKey, waKey, global, local, intervalo)
	return err
}

func ObtenerConfig(db *sql.DB) (phone, newsKey, waKey, global, local string, intervalo int, err error) {
	query := `SELECT phone, news_api_key, wa_api_key, global_topics, local_topics, interval_minutes FROM user_config WHERE id = 1;`
	err = db.QueryRow(query).Scan(&phone, &newsKey, &waKey, &global, &local, &intervalo)
	return
}

func ArticleExists(database *sql.DB, url string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM sent_articles WHERE url=?);`
	err := database.QueryRow(query, url).Scan(&exists)
	if err != nil { return false }
	return exists
}

func MarkAsSent(database *sql.DB, url string) {
	query := `INSERT INTO sent_articles (url) VALUES (?);`
	database.Exec(query, url)
}