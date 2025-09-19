package repository

import (
	"database/sql"
	"movies-csv-import/entity"
)

type MovieRepositoryPostgres struct {
	DB         *sql.DB
	insertStmt *sql.Stmt
}

func NewMovieRepositoryPostgres(db *sql.DB) (*MovieRepositoryPostgres, error) {
	insertStmt, err := db.Prepare("INSERT INTO movies (id, title, year, genres) values ($1, $2, $3, $4)")
	if err != nil {
		return nil, err
	}
	return &MovieRepositoryPostgres{
		DB:         db,
		insertStmt: insertStmt,
	}, nil
}

func (r *MovieRepositoryPostgres) Save(movie entity.Movie) (err error) {
	_, err = r.insertStmt.Exec(movie.ID, movie.Title, movie.Year, movie.Genres)
	return
}

func (r *MovieRepositoryPostgres) Close() error {
	if r.insertStmt != nil {
		return r.insertStmt.Close()
	}
	return nil
}
