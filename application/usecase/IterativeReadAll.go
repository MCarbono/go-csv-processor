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
}

func NewIterativeReadAll(movieRepository repository.MovieRepository) *IterativeReadAll {
	return &IterativeReadAll{
		movieRepository: movieRepository,
	}
}

func (uc *IterativeReadAll) Execute(file *os.File) {
	runtime.GOMAXPROCS(1)
	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	movies := make([]entity.Movie, 0, len(rows))
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
			movies = make([]entity.Movie, 0, 2000)
		}
	}
	if err := uc.movieRepository.SaveBatch(movies); err != nil {
		fmt.Println(err)
	}
}
