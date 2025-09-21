package usecase

import (
	"encoding/csv"
	"fmt"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"os"
	"runtime"
	"sync"
)

type FanOutWorkerReadAll struct {
	movieRepository repository.MovieRepository
	batchSize       int
}

func NewFanOutWorkerReadAll(movieRepository repository.MovieRepository, batchSize int) *FanOutWorkerReadAll {
	return &FanOutWorkerReadAll{
		movieRepository: movieRepository,
		batchSize:       batchSize,
	}
}

func (uc *FanOutWorkerReadAll) Execute(file *os.File) {
	totalWorkers := runtime.GOMAXPROCS(6)
	inDispatcher := make(chan []string, 10*totalWorkers)
	outDispatcher := make(chan entity.Movie, 10*totalWorkers)
	go uc.dispatcher(inDispatcher, outDispatcher)
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

func (uc *FanOutWorkerReadAll) dispatcher(in chan []string, out chan entity.Movie) {
	defer close(out)
	for msg := range in {
		m, err := entity.NewMovie(msg[0], msg[1], msg[2])
		if err != nil {
			fmt.Errorf("error creating movie. %w", err)
			continue
		}
		out <- m
	}
}

func (uc *FanOutWorkerReadAll) worker(in chan entity.Movie, wg *sync.WaitGroup) {
	defer wg.Done()

	batch := make([]entity.Movie, 0, uc.batchSize)

	for msg := range in {
		batch = append(batch, msg)

		// Flush when batch is full
		if len(batch) >= uc.batchSize {
			if err := uc.movieRepository.SaveBatch(batch); err != nil {
				fmt.Printf("Error saving batch: %v\n", err)
			}
			batch = batch[:0] // Reset slice but keep capacity
		}
	}

	// Channel closed, flush remaining batch if there are items
	if len(batch) > 0 {
		if err := uc.movieRepository.SaveBatch(batch); err != nil {
			fmt.Printf("Error saving batch: %v\n", err)
		}
	}
}

// func (uc *FanOutWorkerReadAll) worker(in chan entity.Movie, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for msg := range in {
// 		if err := uc.movieRepository.Save(msg); err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }
