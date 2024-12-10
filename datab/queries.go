package datab

import (
	"database/sql"
	"fmt"
	"strings"

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

// FetchLatestTopRatedMovies retrieves the last 50 (or fewer) movies added with IMDb rating greater than 6
func FetchLatestTopRatedMovies() ([]map[string]interface{}, error) {
	query := `
		SELECT code, title, rating, year, image_link
		FROM movies
		WHERE rating::float > 6
		ORDER BY created_at DESC, rating::float DESC
		LIMIT 50;
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest top-rated movies: %w", err)
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

	// Проверяем, произошли ли ошибки при итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return movies, nil
}

// GetMovieByCode retrieves a movie from the database by its unique code.
func GetMovieByCode(code string) (*Movie, error) {
	query := `SELECT code, title, rating, year, image_link FROM movies WHERE code = $1`
	var movie Movie
	err := db.QueryRow(query, code).Scan(&movie.Code, &movie.Title, &movie.Rating, &movie.Year, &movie.ImageLink)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("movie with code %s not found", code)
		}
		return nil, fmt.Errorf("failed to fetch movie: %w", err)
	}
	return &movie, nil
}

func SearchMovies(query string, years []string, minRating float64) ([]Movie, error) {
	// Строим базовый запрос
	sqlQuery := `
		SELECT code, title, rating, year, image_link
		FROM movies
		WHERE title ILIKE '%' || $1 || '%'
	`

	// Сначала фильтруем по названию
	args := []interface{}{query}

	// Фильтрация по годам
	if len(years) > 0 {
		sqlQuery += " AND year = ANY($2::text[])"
		yearArray := "{" + strings.Join(years, ",") + "}"
		args = append(args, yearArray)
	}

	// Фильтрация по рейтингу
	if minRating > 0 {
		sqlQuery += " AND rating::float >= $3"
		args = append(args, minRating)
	}

	// Выполняем запрос
	rows, err := db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var movie Movie
		if err := rows.Scan(&movie.Code, &movie.Title, &movie.Rating, &movie.Year, &movie.ImageLink); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		movies = append(movies, movie)
	}

	// Проверка на ошибки при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return movies, nil
}
