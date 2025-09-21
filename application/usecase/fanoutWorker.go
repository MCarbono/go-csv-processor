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
	totalWorkers := runtime.NumCPU() * 2
	inDispatcher := make(chan []string, 10*totalWorkers)
	outDispatcher := make(chan entity.Movie, 10*totalWorkers)
	go uc.dispatcher(inDispatcher, outDispatcher)
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

func (uc *FanOutWorker) dispatcher(in chan []string, out chan entity.Movie) {
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

func (uc *FanOutWorker) worker(in chan entity.Movie, wg *sync.WaitGroup) {
	defer wg.Done()
	batch := make([]entity.Movie, 0, 1000)

	for msg := range in {
		batch = append(batch, msg)

		// Flush when batch is full
		if len(batch) >= 1000 {
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

// func (uc *FanOutWorker) worker(in chan entity.Movie, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for msg := range in {
// 		if err := uc.movieRepository.Save(msg); err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }
