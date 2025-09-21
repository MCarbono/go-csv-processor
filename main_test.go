package main

import (
	"fmt"
	"movies-csv-import/entity"
	"movies-csv-import/infra/database"
	"movies-csv-import/infra/repository"
	"testing"
)

// import (
// 	"movies-csv-import/application/factory"
// 	"movies-csv-import/infra/database"
// 	"movies-csv-import/infra/repository"
// 	"os"
// 	"testing"
// )

// func BenchmarkIterative(b *testing.B) {
// 	DB, err := database.Open()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer DB.Close()
// 	f, err := os.Open("movie.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	movieRepository := repository.NewMovieRepositoryPostgres(DB)
// 	uc, err := factory.CreateUseCase(movieRepository, factory.ITERATIVE)
// 	if err != nil {
// 		panic(err)
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		uc.Execute(f)
// 	}
// }

// func BenchmarkIterativeReadAll(b *testing.B) {
// 	DB, err := database.Open()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer DB.Close()
// 	f, err := os.Open("movie.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	movieRepository := repository.NewMovieRepositoryPostgres(DB)
// 	uc, err := factory.CreateUseCase(movieRepository, factory.ITERATIVE_READALL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < b.N; i++ {
// 		uc.Execute(f)
// 	}
// 	b.ResetTimer()
// }

// func BenchmarkIPipelineWorker(b *testing.B) {
// 	DB, err := database.Open()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer DB.Close()
// 	f, err := os.Open("movie.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	movieRepository := repository.NewMovieRepositoryPostgres(DB)
// 	uc, err := factory.CreateUseCase(movieRepository, factory.PIPELINE_WORKER)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for i := 0; i < b.N; i++ {
// 		uc.Execute(f)
// 	}
// 	b.ResetTimer()
// }

// func BenchmarkIPipelineWorkerReadAll(b *testing.B) {
// 	DB, err := database.Open()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer DB.Close()
// 	f, err := os.Open("movie.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	movieRepository := repository.NewMovieRepositoryPostgres(DB)
// 	uc, err := factory.CreateUseCase(movieRepository, factory.PIPELINE_WORKER_READALL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		uc.Execute(f)
// 	}
// }

// func BenchmarkFanOutWorker(b *testing.B) {
// 	DB, err := database.Open()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer DB.Close()
// 	f, err := os.Open("movie.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	movieRepository := repository.NewMovieRepositoryPostgres(DB)
// 	uc, err := factory.CreateUseCase(movieRepository, factory.FANOUT_WORKER)
// 	if err != nil {
// 		panic(err)
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		uc.Execute(f)
// 	}
// }

// func BenchmarkFanOutWorkerReadAll(b *testing.B) {
// 	DB, err := database.Open()
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer DB.Close()
// 	f, err := os.Open("movie.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f.Close()
// 	movieRepository := repository.NewMovieRepositoryPostgres(DB)
// 	uc, err := factory.CreateUseCase(movieRepository, factory.FANOUT_WORKER_READALL)
// 	if err != nil {
// 		panic(err)
// 	}
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		uc.Execute(f)
// 	}
// }

// Add this to your main_test.go or create a new benchmark file

func BenchmarkRepositoryWithPreparedStmt(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	// Your current optimized implementation
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB)
	if err != nil {
		panic(err)
	}
	defer movieRepository.Close()

	movie := entity.Movie{
		ID:     "1",
		Title:  "Test Movie (2023)",
		Year:   "2023",
		Genres: "Action",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		movie.ID = fmt.Sprintf("%d", i)
		_ = movieRepository.Save(movie)
	}
}

func BenchmarkRepositoryWithoutPreparedStmt(b *testing.B) {
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}
	defer DB.Close()

	// Simulate the old implementation
	movie := entity.Movie{
		ID:     "1",
		Title:  "Test Movie (2023)",
		Year:   "2023",
		Genres: "Action",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		movie.ID = fmt.Sprintf("%d", i)
		// This compiles the query every time (inefficient)
		_, _ = DB.Exec("INSERT INTO movies (id, title, year, genres) VALUES ($1, $2, $3, $4)",
			movie.ID, movie.Title, movie.Year, movie.Genres)
	}
}
