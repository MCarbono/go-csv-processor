package repository

import (
	"database/sql"
	"movies-csv-import/entity"
)

type MovieRepositoryPostgres struct {
	DB *sql.DB
}

func NewMovieRepositoryPostgres(db *sql.DB) *MovieRepositoryPostgres {
	return &MovieRepositoryPostgres{
		DB: db,
	}
}

func (r *MovieRepositoryPostgres) Save(movie *entity.Movie) (err error) {
	_, err = r.DB.Exec("INSERT INTO movies (id, title, year, genres) values ($1, $2, $3, $4)", movie.ID, movie.Title, movie.Year, movie.Genres)
	return
}
