<h1 align="center">CSV Movies</h1>


## 📃 About

This is a basic project to practice GO concurrency/paralellism. The challenge was to read from a CSV file
several movies and storing each one of them in a database. I developed 6 possible ways of doing it, 4 of them using 
concurrency. This code is a challenge of the https://app.devgym.com.br/ platform


## 🗄 Libs/Dependencies

| Name        | Description | Documentation | Installation |
| ----------- | ----------- | ------------- | ----------- |     
| pgx      | postgres database driver       |  github.com/jackc/pgx/v4 |  go get go get github.com/jackc/pgx/v4      |

## ⚙️ Setup

Clone the github repository:

```bash
    $ git clone https://github.com/MCarbono/go-csv-processor.git
``` 

Go to project's folder

```bash
    $ cd go-csv-processor
```

Start the database with one of the commands below: 

```bash
    # docker command
    $ docker compose up -d
```

```bash
    # Makefile comando
    $ make db_up
```

To exclude the database container, use one of the commands below: 

```bash
    # docker command
    $ docker compose down
```

```bash
    # Makefile command
    $ make db_down
```

## ⚙️ Run

Use one of the commands below to run one usecase: 

```bash
    # Makefile command
    $ make iterative
```

```bash
    # Makefile command
    $ make iterative-readall
```

```bash
    # Makefile command
    $ make pipeline-worker-readall
```

```bash
    # Makefile command
    $ make pipeline-worker
```

```bash
    # Makefile command
    $ make fanout-worker
```

```bash
    # Makefile command
    $ make fanout-worker-readall
```

## 🧪 Benchmark

In the root project folder, run one of the commands below:

```bash
    # Makefile command
    $ make bench
```

```bash
    # Go command
    $ go test -bench=. -benchmem
```

