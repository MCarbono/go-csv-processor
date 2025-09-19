package entity

import (
	"regexp"
	"strings"
)

type Movie struct {
	ID     string
	Title  string
	Year   string
	Genres string
}

var findYearRegex = regexp.MustCompile(`\((\d+)\)`)

func NewMovie(ID, title, genres string) (Movie, error) {
	m := Movie{
		ID:     ID,
		Title:  title,
		Genres: genres,
	}
	if matches := findYearRegex.FindStringSubmatch(title); len(matches) > 1 {
		m.Year = matches[1]
		m.Title = strings.TrimSpace(findYearRegex.ReplaceAllString(title, ""))
	}
	return m, nil
}
