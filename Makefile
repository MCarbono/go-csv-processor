# Build target
build:
	go build -o movies-csv-import main.go

# Clean target
clean:
	rm -f movies-csv-import

# Development targets (using go run)
iterative:
	go run main.go

iterative-readall:
	go run main.go --usecase=iterative-readall
	
pipeline-worker-readall:
	go run main.go --usecase=pipeline-worker-readall

pipeline-worker:
	go run main.go --usecase=pipeline-worker-streaming

fanout-worker:
	go run main.go --usecase=fanout-worker

fanout-worker-readall:
	go run main.go --usecase=fanout-worker-readall

# Production targets (using built binary)
run-iterative: build
	time ./movies-csv-import

run-iterative-readall: build
	time ./movies-csv-import --usecase=iterative-readall

run-pipeline-worker-readall: build
	time ./movies-csv-import --usecase=pipeline-worker-readall

run-pipeline-worker: build
	time ./movies-csv-import --usecase=pipeline-worker-streaming

run-fanout-worker: build
	time ./movies-csv-import --usecase=fanout-worker

run-fanout-worker-readall: build
	time ./movies-csv-import --usecase=fanout-worker-readall

db_down:
	docker compose down 

db_up:
	docker compose up -d

bench:
	go test -bench=. -benchmem

bench-profile:
	go test -bench=. -benchmem -cpuprofile=cpu.prof -memprofile=mem.prof

bench-race:
	go test -bench=. -race

deps:
	go mod download
	go mod tidy