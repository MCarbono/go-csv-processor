package main

import (
	"flag"
	"fmt"
	"log"
	"math/big"
	"movies-csv-import/application/factory"
	"movies-csv-import/infra/database"
	"movies-csv-import/infra/repository"
	"os"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	usecaseType = flag.String("usecase", "iterative", "which usecase should use")
)

func main() {
	start := time.Now()
	flag.Parse()
	fmt.Printf("import method choosed is: %v\n", *usecaseType)
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
	uc, err := factory.CreateUseCase(movieRepository, *usecaseType)
	if err != nil {
		panic(err)
	}
	uc.Execute(f)
	r := new(big.Int)
	fmt.Println(r.Binomial(1000, 10))
	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
}

//total com os erros que precisam estar inseridos no banco: 27275
//ids dos registros com errors:
//7789
//51372
//112809
