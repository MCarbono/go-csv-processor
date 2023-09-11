package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"os"
)

type Iterative struct {
	movieRepository repository.MovieRepository
}

func NewIterative(movieRepository repository.MovieRepository) *Iterative {
	return &Iterative{
		movieRepository: movieRepository,
	}
}

func (uc *Iterative) Execute(file *os.File) {
	rows := csv.NewReader(file)
	//read first line - header
	_, err := rows.Read()
	if err != nil {
		if err != io.EOF {
			log.Fatalln(err)
		}
	}
	for {
		row, err := rows.Read()
		if err == io.EOF {
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
		if err := uc.movieRepository.Save(movie); err != nil {
			fmt.Println(err)
		}
	}
}
