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

func NewMovie(ID, title, genres string) (*Movie, error) {
	m := &Movie{
		ID:     ID,
		Title:  title,
		Genres: genres,
	}
	r, err := regexp.Compile(`\(\d*\)`)
	if err != nil {
		return nil, err
	}
	year := r.FindString(title)
	if year != "" {
		m.Title = r.ReplaceAllString(m.Title, "")
		m.Title = strings.TrimSpace(m.Title)
		year = strings.Replace(year, "(", "", 1)
		year = strings.Replace(year, ")", "", 1)
		m.Year = year
	}
	return m, nil
}
