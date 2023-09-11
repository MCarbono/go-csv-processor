package main

import (
	"movies-csv-import/application/factory"
	"movies-csv-import/infra/database"
	"movies-csv-import/infra/repository"
	"os"
	"testing"
)

func BenchmarkIterative(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	f, err := os.Open("movie.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	movieRepository := repository.NewMovieRepositoryPostgres(DB)
	uc, err := factory.CreateUseCase(movieRepository, factory.ITERATIVE)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
	}
}

func BenchmarkIterativeReadAll(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	f, err := os.Open("movie.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	movieRepository := repository.NewMovieRepositoryPostgres(DB)
	uc, err := factory.CreateUseCase(movieRepository, factory.ITERATIVE_READALL)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
	}
}

func BenchmarkIPipelineWorker(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	f, err := os.Open("movie.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	movieRepository := repository.NewMovieRepositoryPostgres(DB)
	uc, err := factory.CreateUseCase(movieRepository, factory.PIPELINE_WORKER)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
	}
}

func BenchmarkIPipelineWorkerReadAll(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	f, err := os.Open("movie.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	movieRepository := repository.NewMovieRepositoryPostgres(DB)
	uc, err := factory.CreateUseCase(movieRepository, factory.PIPELINE_WORKER_READALL)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
	}
}

func BenchmarkFanOutWorker(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	f, err := os.Open("movie.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	movieRepository := repository.NewMovieRepositoryPostgres(DB)
	uc, err := factory.CreateUseCase(movieRepository, factory.FANOUT_WORKER)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
	}
}

func BenchmarkFanOutWorkerReadAll(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	f, err := os.Open("movie.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	movieRepository := repository.NewMovieRepositoryPostgres(DB)
	uc, err := factory.CreateUseCase(movieRepository, factory.FANOUT_WORKER_READALL)
	if err != nil {
		panic(err)
	}
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
	}
}
