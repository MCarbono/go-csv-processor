package usecase

import (
	"fmt"
	"movies-csv-import/application/repository"
	"movies-csv-import/entity"
	"sync"
)

type WorkerLauncher interface {
	LaunchWorker(in chan []string)
}

type worker struct {
	movieRepository repository.MovieRepository
	wg              *sync.WaitGroup
	batchSize       int
}

func NewWorker(movieRepository repository.MovieRepository, wg *sync.WaitGroup, batchSize int) WorkerLauncher {
	return &worker{movieRepository: movieRepository, wg: wg, batchSize: batchSize}
}

func (w *worker) LaunchWorker(in chan []string) {
	w.Save(w.DataTransform(in))
}

func (w *worker) DataTransform(in <-chan []string) chan entity.Movie {
	out := make(chan entity.Movie)
	go func() {
		defer close(out)
		for msg := range in {
			m, err := entity.NewMovie(msg[0], msg[1], msg[2])
			if err != nil {
				fmt.Errorf("error creating movie. %w", err)
				continue
			}
			out <- m
		}
	}()
	return out
}

func (w *worker) Save(in chan entity.Movie) {
	go func() {
		defer w.wg.Done()

		batch := make([]entity.Movie, 0, w.batchSize)

		for msg := range in {
			batch = append(batch, msg)

			// Flush when batch is full
			if len(batch) >= w.batchSize {
				if err := w.movieRepository.SaveBatch(batch); err != nil {
					fmt.Printf("Error saving batch: %v\n", err)
				}
				batch = batch[:0] // Reset slice but keep capacity
			}
		}

		// Channel closed, flush remaining batch if there are items
		if len(batch) > 0 {
			if err := w.movieRepository.SaveBatch(batch); err != nil {
				fmt.Printf("Error saving batch: %v\n", err)
			}
		}
	}()
}

// func (w *worker) Save(in <-chan entity.Movie) {
// 	go func() {
// 		defer w.wg.Done()
// 		for msg := range in {
// 			if err := w.movieRepository.Save(msg); err != nil {
// 				fmt.Println(err)
// 			}
// 		}
// 	}()
// }
