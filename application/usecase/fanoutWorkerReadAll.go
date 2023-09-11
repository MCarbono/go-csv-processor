package usecase

import (
	"encoding/csv"
	"fmt"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"os"
	"sync"
)

type FanOutWorkerReadAll struct {
	movieRepository repository.MovieRepository
}

func NewFanOutWorkerReadAll(movieRepository repository.MovieRepository) *FanOutWorkerReadAll {
	return &FanOutWorkerReadAll{
		movieRepository: movieRepository,
	}
}

func (uc *FanOutWorkerReadAll) Execute(file *os.File) {
	inDispatcher := make(chan []string, 100)
	outDispatcher := make(chan *entity.Movie, 100)
	go uc.dispatcher(inDispatcher, outDispatcher)
	totalWorkers := 18
	var wg sync.WaitGroup
	wg.Add(totalWorkers)
	for i := 0; i < totalWorkers; i++ {
		go uc.worker(outDispatcher, &wg)
	}
	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(rows); i++ {
		inDispatcher <- rows[i]
	}
	close(inDispatcher)
	wg.Wait()
}

func (uc *FanOutWorkerReadAll) dispatcher(in chan []string, out chan *entity.Movie) {
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

func (uc *FanOutWorkerReadAll) worker(in chan *entity.Movie, wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range in {
		if err := uc.movieRepository.Save(msg); err != nil {
			fmt.Println(err)
		}
	}
}
