package usecase

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"os"
	"sync"
)

type FanOutWorker struct {
	movieRepository repository.MovieRepository
}

func NewFanOutWorker(movieRepository repository.MovieRepository) *FanOutWorker {
	return &FanOutWorker{
		movieRepository: movieRepository,
	}
}

func (uc *FanOutWorker) Execute(file *os.File) {
	inDispatcher := make(chan []string, 100)
	outDispatcher := make(chan *entity.Movie, 100)
	go uc.dispatcher(inDispatcher, outDispatcher)
	totalWorkers := 18
	var wg sync.WaitGroup
	wg.Add(totalWorkers)
	for i := 0; i < totalWorkers; i++ {
		go uc.worker(outDispatcher, &wg)
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
		inDispatcher <- row
	}
	close(inDispatcher)
	wg.Wait()

}

func (uc *FanOutWorker) dispatcher(in chan []string, out chan *entity.Movie) {
	defer close(out)
	for msg := range in {
		m, err := entity.NewMovie(string(msg[0]), string(msg[1]), string(msg[2]))
		if err != nil {
			fmt.Errorf("error creating movie. %w", err)
			continue
		}
		out <- m
	}
}

func (uc *FanOutWorker) worker(in chan *entity.Movie, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range in {
		if err := uc.movieRepository.Save(msg); err != nil {
			fmt.Println(err)
		}
	}
}
