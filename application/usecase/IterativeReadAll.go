package usecase

import (
	"encoding/csv"
	"fmt"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"os"
	"runtime"
)

type IterativeReadAll struct {
	movieRepository repository.MovieRepository
	batchSize       int
}

func NewIterativeReadAll(movieRepository repository.MovieRepository, batchSize int) *IterativeReadAll {
	return &IterativeReadAll{
		movieRepository: movieRepository,
		batchSize:       batchSize,
	}
}

func (uc *IterativeReadAll) Execute(file *os.File) {
	runtime.GOMAXPROCS(6)
	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	movies := make([]entity.Movie, 0, uc.batchSize)
	for i := 1; i < len(rows); i++ {
		movie, err := entity.NewMovie(rows[i][0], rows[i][1], rows[i][2])
		if err != nil {
			fmt.Errorf("error creating movie. %w", err)
			continue
		}
		movies = append(movies, movie)
		if len(movies) >= 2000 {
			if err := uc.movieRepository.SaveBatch(movies); err != nil {
				fmt.Println(err)
			}
			movies = movies[:0]
		}
	}
	if err := uc.movieRepository.SaveBatch(movies); err != nil {
		fmt.Println(err)
	}
}
