package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"movies-csv-import/application/repository"
	"os"
	"sync"
)

type PipelineWorker struct {
	movieRepository repository.MovieRepository
}

func NewPipelineWorker(movieRepository repository.MovieRepository) *PipelineWorker {
	return &PipelineWorker{
		movieRepository: movieRepository,
	}
}

func (uc *PipelineWorker) Execute(file *os.File) {
	dispatcher := NewDispatcher(100)
	totalWorkers := 18
	var wg sync.WaitGroup
	wg.Add(totalWorkers)
	for i := 0; i < totalWorkers; i++ {
		worker := NewWorker(uc.movieRepository, &wg)
		dispatcher.LaunchWorker(worker)
	}
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
		dispatcher.Launch(row)
	}
	dispatcher.Stop()
	wg.Wait()
}
