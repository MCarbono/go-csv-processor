package usecase

import "os"

type CSVFile interface {
	Execute(file *os.File)
}
