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
}

func NewWorker(movieRepository repository.MovieRepository, wg *sync.WaitGroup) WorkerLauncher {
	return &worker{movieRepository: movieRepository, wg: wg}
}

func (w *worker) LaunchWorker(in chan []string) {
	w.Save(w.DataTransform(in))
}

func (w *worker) DataTransform(in <-chan []string) chan *entity.Movie {
	out := make(chan *entity.Movie)
	go func() {
		defer close(out)
		for msg := range in {
			m, err := entity.NewMovie(string(msg[0]), string(msg[1]), string(msg[2]))
			if err != nil {
				fmt.Errorf("error creating movie. %w", err)
				continue
			}
			out <- m
		}
	}()
	return out
}

func (w *worker) Save(in <-chan *entity.Movie) {
	go func() {
		defer w.wg.Done()
		for msg := range in {
			if err := w.movieRepository.Save(msg); err != nil {
				fmt.Println(err)
			}
		}
	}()
}
