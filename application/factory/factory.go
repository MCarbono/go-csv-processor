package factory

import (
	"fmt"
	"movies-csv-import/application/repository"
	"movies-csv-import/application/usecase"
)

const (
	ITERATIVE               = "iterative"
	ITERATIVE_READALL       = "iterative-readall"
	PIPELINE_WORKER_READALL = "pipeline-worker-readall"
	PIPELINE_WORKER         = "pipeline-worker-streaming"
	FANOUT_WORKER           = "fanout-worker"
	FANOUT_WORKER_READALL   = "fanout-worker-readall"
)

func CreateUseCase(movieRepository repository.MovieRepository, usecaseType string, batchSize int) (usecase.CSVFile, error) {
	if usecaseType == PIPELINE_WORKER_READALL {
		return usecase.NewPipelineWorkerReadAll(movieRepository, batchSize), nil
	}
	if usecaseType == PIPELINE_WORKER {
		return usecase.NewPipelineWorker(movieRepository, batchSize), nil
	}
	if usecaseType == FANOUT_WORKER {
		return usecase.NewFanOutWorker(movieRepository, batchSize), nil
	}
	if usecaseType == FANOUT_WORKER_READALL {
		return usecase.NewFanOutWorkerReadAll(movieRepository, batchSize), nil
	}
	if usecaseType == ITERATIVE_READALL {
		return usecase.NewIterativeReadAll(movieRepository, batchSize), nil
	}
	if usecaseType == ITERATIVE {
		return usecase.NewIterative(movieRepository, batchSize), nil
	}
	return nil, fmt.Errorf("usecase type %v not found", usecaseType)
}
