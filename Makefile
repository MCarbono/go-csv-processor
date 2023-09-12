
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

db_down:
	docker compose down 

db_up:
	docker compose up -d

bench:
	go test -bench=. -benchmem