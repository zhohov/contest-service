package domain

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	maxProblemTitleLength = 256
	maxProblemSlugLength  = 256
)

type Problem struct {
	ID         int
	Title      string
	Slug       string
	Difficulty string
	Content    string
	TestCases  []TestCase
}

func NewProblem(
	title, slug, difficulty, content string,
	testCases []TestCase,
) (*Problem, error) {
	problem := &Problem{
		Title:      title,
		Slug:       slug,
		Difficulty: difficulty,
		Content:    content,
		TestCases:  testCases,
	}

	if err := problem.validate(); err != nil {
		return nil, err
	}

	return problem, nil
}

func (p *Problem) validate() error {
	if strings.TrimSpace(p.Title) == "" {
		return errors.New("problem title is required")
	}

	if strings.TrimSpace(p.Slug) == "" {
		return errors.New("problem slug is required")
	}

	if utf8.RuneCountInString(p.Title) > maxProblemTitleLength {
		return fmt.Errorf("problem title length must be <= %d", maxProblemTitleLength)
	}

	if utf8.RuneCountInString(p.Slug) > maxProblemSlugLength {
		return fmt.Errorf("problem slug length must be <= %d", maxProblemSlugLength)
	}

	return nil
}
