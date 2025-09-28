package main

import (
	"movies-csv-import/application/factory"
	"movies-csv-import/infra/database"
	"movies-csv-import/infra/repository"
	"os"
	"testing"
)

func BenchmarkIterative(b *testing.B) {
	batchSize := 2000
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
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB, batchSize)
	if err != nil {
		panic(err)
	}
	uc, err := factory.CreateUseCase(movieRepository, factory.ITERATIVE, batchSize)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
		_, err = DB.Exec("DELETE FROM movies;")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkIterativeReadAll(b *testing.B) {
	batchSize := 2000
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
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB, batchSize)
	if err != nil {
		panic(err)
	}
	uc, err := factory.CreateUseCase(movieRepository, factory.ITERATIVE_READALL, batchSize)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
		_, err = DB.Exec("DELETE FROM movies;")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkIPipelineWorker(b *testing.B) {
	batchSize := 2000
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
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB, batchSize)
	if err != nil {
		panic(err)
	}
	uc, err := factory.CreateUseCase(movieRepository, factory.PIPELINE_WORKER, batchSize)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
		_, err = DB.Exec("DELETE FROM movies;")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkIPipelineWorkerReadAll(b *testing.B) {
	batchSize := 2000
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
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB, batchSize)
	if err != nil {
		panic(err)
	}
	uc, err := factory.CreateUseCase(movieRepository, factory.PIPELINE_WORKER_READALL, batchSize)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
		_, err = DB.Exec("DELETE FROM movies;")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFanOutWorker(b *testing.B) {
	batchSize := 2000
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
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB, batchSize)
	if err != nil {
		panic(err)
	}
	uc, err := factory.CreateUseCase(movieRepository, factory.FANOUT_WORKER, batchSize)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
		_, err = DB.Exec("DELETE FROM movies;")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkFanOutWorkerReadAll(b *testing.B) {
	batchSize := 2000
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
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB, batchSize)
	if err != nil {
		panic(err)
	}
	uc, err := factory.CreateUseCase(movieRepository, factory.FANOUT_WORKER_READALL, batchSize)
	if err != nil {
		panic(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		uc.Execute(f)
		_, err = DB.Exec("DELETE FROM movies;")
		if err != nil {
			panic(err)
		}
	}
}
