package usecase

import (
	"encoding/csv"
	"fmt"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"os"
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
	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(rows); i++ {
		movie, err := entity.NewMovie(rows[i][0], rows[i][1], rows[i][2])
		if err != nil {
			fmt.Errorf("error creating movie. %w", err)
			continue
		}
		if err := uc.movieRepository.Save(movie); err != nil {
			fmt.Println(err)
		}
	}
}
