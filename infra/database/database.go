package database

import "database/sql"

func Open() (DB *sql.DB, err error) {
	DB, err = sql.Open("pgx", "host=localhost port=5432 user=movies password=movies dbname=movies sslmode=disable")
	if err != nil {
		return
	}
	_, err = DB.Exec("DELETE FROM movies;")
	if err != nil {
		return
	}
	return
}
