package domain

import "time"

type Contest struct {
	ID          int
	Creator     User
	Title       string
	Slug        string
	Description string
	Problems    []Problem
	StartsAt    time.Time
	EndsAt      time.Time
	CreatedAt   time.Time
}

func NewContest(
	creator User,
	title string,
	slug string,
	description string,
	problems []Problem,
	startsAt time.Time,
	endsAt time.Time,
) *Contest {
	return &Contest{
		Creator:     creator,
		Title:       title,
		Slug:        slug,
		Description: description,
		Problems:    problems,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
	}
}

func (c *Contest) IsActive(now time.Time) bool {
	return now.After(c.StartsAt) && now.Before(c.EndsAt)
}
