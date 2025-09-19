package repository

import "movies-csv-import/entity"

type MovieRepository interface {
	Save(movie entity.Movie) (err error)
}
