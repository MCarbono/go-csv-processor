package repository

import (
	"database/sql"
	"fmt"
	"movies-csv-import/entity"
	"strings"
)

type MovieRepositoryPostgres struct {
	DB              *sql.DB
	insertStmt      *sql.Stmt
	batchInsertStmt *sql.Stmt
	batchSize       int
}

func NewMovieRepositoryPostgres(db *sql.DB, batchSize int) (*MovieRepositoryPostgres, error) {
	insertStmt, err := db.Prepare("INSERT INTO movies (id, title, year, genres) values ($1, $2, $3, $4)")
	if err != nil {
		return nil, err
	}
	// Pre-prepare batch insert statement for maximum batch size
	batchInsertStmt, err := db.Prepare(generateBatchQuery(1000))
	if err != nil {
		return nil, err
	}

	return &MovieRepositoryPostgres{
		DB:              db,
		insertStmt:      insertStmt,
		batchInsertStmt: batchInsertStmt,
	}, nil
}

func generateBatchQuery(count int) string {
	valueStrings := make([]string, count)
	for i := 0; i < count; i++ {
		valueStrings[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4)
	}
	return "INSERT INTO movies (id, title, year, genres) VALUES " + strings.Join(valueStrings, ",")
}

func (r *MovieRepositoryPostgres) Save(movie entity.Movie) (err error) {
	_, err = r.insertStmt.Exec(movie.ID, movie.Title, movie.Year, movie.Genres)
	return
}

func (r *MovieRepositoryPostgres) Close() error {
	var err error
	if r.insertStmt != nil {
		err = r.insertStmt.Close()
	}
	if r.batchInsertStmt != nil {
		if closeErr := r.batchInsertStmt.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}

// Optimized batch insert with prepared statement
func (r *MovieRepositoryPostgres) SaveBatch(movies []entity.Movie) (err error) {
	if len(movies) == 0 {
		return nil
	}

	// Use prepared statement for common batch sizes
	if len(movies) == r.batchSize {
		valueArgs := make([]any, 0, len(movies)*4)
		for _, movie := range movies {
			valueArgs = append(valueArgs, movie.ID, movie.Title, movie.Year, movie.Genres)
		}

		tx, err := r.DB.Begin()
		if err != nil {
			return err
		}

		_, err = tx.Stmt(r.batchInsertStmt).Exec(valueArgs...)
		if err != nil {
			tx.Rollback()
			return err
		}

		return tx.Commit()
	}

	// Fallback for other batch sizes
	return r.SaveBatchDynamic(movies)
}

func (r *MovieRepositoryPostgres) SaveBatchDynamic(movies []entity.Movie) (err error) {
	valueStrings := make([]string, 0, len(movies))
	valueArgs := make([]any, 0, len(movies)*4)

	for i, movie := range movies {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)",
			i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, movie.ID, movie.Title, movie.Year, movie.Genres)
	}

	query := "INSERT INTO movies (id, title, year, genres) VALUES " + strings.Join(valueStrings, ",")

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				fmt.Printf("error rolling back transaction: %v\n", errRollback)
			}
		}
	}()

	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return tx.Commit()
}
