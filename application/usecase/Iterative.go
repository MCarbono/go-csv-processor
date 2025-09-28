package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"os"
	"runtime"
)

type Iterative struct {
	movieRepository repository.MovieRepository
	batchSize       int
}

func NewIterative(movieRepository repository.MovieRepository, batchSize int) *Iterative {
	return &Iterative{
		movieRepository: movieRepository,
		batchSize:       batchSize,
	}
}

func (uc *Iterative) Execute(file *os.File) {
	runtime.GOMAXPROCS(6)
	rows := csv.NewReader(file)
	//read first line - header
	_, err := rows.Read()
	if err != nil {
		if err != io.EOF {
			log.Fatalln(err)
		}
	}
	movies := make([]entity.Movie, 0, uc.batchSize)
	for {
		row, err := rows.Read()
		if err == io.EOF {
			if len(movies) > 0 {
				if err := uc.movieRepository.SaveBatch(movies); err != nil {
					fmt.Println(err)
				}
			}
			break
		}
		if err != nil {
			fmt.Println(err)
			continue
		}

		movie, err := entity.NewMovie(row[0], row[1], row[2])
		if err != nil {
			fmt.Errorf("error creating movie. %w", err)
			continue
		}
		movies = append(movies, movie)
		if len(movies) >= uc.batchSize {
			if err := uc.movieRepository.SaveBatch(movies); err != nil {
				fmt.Println(err)
			}
			movies = movies[:0]
		}
	}
}
