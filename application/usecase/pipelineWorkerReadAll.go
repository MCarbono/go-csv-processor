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
	batchSize       int
}

func NewPipelineWorkerReadAll(movieRepository repository.MovieRepository, batchSize int) *PipelineWorkerReadlAll {
	return &PipelineWorkerReadlAll{
		movieRepository: movieRepository,
		batchSize:       batchSize,
	}
}

func (uc *PipelineWorkerReadlAll) Execute(file *os.File) {
	// totalWorkers := runtime.NumCPU() * 2
	totalWorkers := runtime.GOMAXPROCS(6)
	dispatcher := NewDispatcher(uc.batchSize)
	var wg sync.WaitGroup
	wg.Add(totalWorkers)
	for i := 0; i < totalWorkers; i++ {
		worker := NewWorker(uc.movieRepository, &wg, uc.batchSize)
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
