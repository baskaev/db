package datab

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() error {
	var err error
	connStr := "host=db user=user password=password dbname=films_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка подключения
	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	return nil
}

// FetchMovies retrieves all movies from the database
func FetchMovies() ([]map[string]interface{}, error) {
	rows, err := db.Query("SELECT id, title, description, release_date FROM movies")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}
	defer rows.Close()

	var movies []map[string]interface{}
	for rows.Next() {
		var id int
		var title, description string
		var releaseDate string
		if err := rows.Scan(&id, &title, &description, &releaseDate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		movies = append(movies, map[string]interface{}{
			"id":           id,
			"title":        title,
			"description":  description,
			"release_date": releaseDate,
		})
	}
	return movies, nil
}
