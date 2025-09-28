package main

import (
	"flag"
	"fmt"
	"movies-csv-import/application/factory"
	"movies-csv-import/infra/database"
	"movies-csv-import/infra/repository"
	"os"
	"runtime/trace"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	usecaseType = flag.String("usecase", "iterative", "which usecase should use")
	traceFile   = flag.String("trace", "trace.out", "trace file to write to")
)

func main() {
	flag.Parse()

	f, err := os.Create(*traceFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer fmt.Printf("import method choosed is: %v\n", *usecaseType)
	DB, err := database.Open()
	if err != nil {
		panic(err)
	}

	f, err = os.Open("movie.csv")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Printf("error closing file: %v\n", err)
		}
		err = DB.Close()
		if err != nil {
			fmt.Printf("error closing database: %v\n", err)
		}
		trace.Stop()
	}()
	batchSize := 1000
	movieRepository, err := repository.NewMovieRepositoryPostgres(DB, batchSize)
	if err != nil {
		panic(err)
	}
	defer movieRepository.Close()
	uc, err := factory.CreateUseCase(movieRepository, *usecaseType, batchSize)
	if err != nil {
		panic(err)
	}
	uc.Execute(f)
}
