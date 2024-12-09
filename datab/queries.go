package datab

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var db *sql.DB

type Movie struct {
	Title     string
	Code      string
	Rating    string
	Year      string
	ImageLink string
}

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
	rows, err := db.Query("SELECT code, title, rating, year, image_link FROM movies")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch movies: %w", err)
	}
	defer rows.Close()

	var movies []map[string]interface{}
	for rows.Next() {
		var code, title, rating, year, imageLink string
		if err := rows.Scan(&code, &title, &rating, &year, &imageLink); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		movies = append(movies, map[string]interface{}{
			"code":       code,
			"title":      title,
			"rating":     rating,
			"year":       year,
			"image_link": imageLink,
		})
	}
	return movies, nil
}

// AddMovie inserts a new movie into the database
func AddMovie(movie Movie) error {
	// SQL query to insert a new movie
	query := `INSERT INTO movies (code, title, rating, year, image_link) 
			  VALUES ($1, $2, $3, $4, $5)`

	// Execute the query with the movie data
	_, err := db.Exec(query, movie.Code, movie.Title, movie.Rating, movie.Year, movie.ImageLink)
	if err != nil {
		return fmt.Errorf("failed to insert movie: %w", err)
	}

	return nil
}
