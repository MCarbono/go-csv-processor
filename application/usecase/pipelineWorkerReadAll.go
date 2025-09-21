package usecase

import (
	"encoding/csv"
	"movies-csv-import/application/repository"
	"os"
	"runtime"
	"sync"
)

type PipelineWorkerReadlAll struct {
	movieRepository repository.MovieRepository
}

func NewPipelineWorkerReadAll(movieRepository repository.MovieRepository) *PipelineWorkerReadlAll {
	return &PipelineWorkerReadlAll{
		movieRepository: movieRepository,
	}
}

func (uc *PipelineWorkerReadlAll) Execute(file *os.File) {
	totalWorkers := runtime.GOMAXPROCS(6)
	// totalWorkers := runtime.NumCPU() * 2
	dispatcher := NewDispatcher(10 * totalWorkers)
	var wg sync.WaitGroup
	wg.Add(totalWorkers)
	for i := 0; i < totalWorkers; i++ {
		worker := NewWorker(uc.movieRepository, &wg)
		dispatcher.LaunchWorker(worker)
	}
	csvReader := csv.NewReader(file)
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic(err)
	}
	for i := 1; i < len(rows); i++ {
		dispatcher.Launch(rows[i])
	}
	dispatcher.Stop()
	wg.Wait()
}
