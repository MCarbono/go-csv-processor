package repository

import "movies-csv-import/entity"

type MovieRepository interface {
	Save(movie entity.Movie) (err error)
	SaveBatch(movies []entity.Movie) (err error)
	SaveBatchDynamic(movies []entity.Movie) (err error)
}
