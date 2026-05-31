package repository

import "github.com/zhohov/contest-service/backend/internal/domain"

type ProblemInMemoryRepo struct {
	storage []domain.Problem
}

func NewProblemInMemoryRepo(storage []domain.Problem) *ProblemInMemoryRepo {
	return &ProblemInMemoryRepo{
		storage: storage,
	}
}

func (r *ProblemInMemoryRepo) CreateProblem(problem *domain.Problem) (*domain.Problem, error) {
	r.storage = append(r.storage, *problem)
	return problem, nil
}

func (r *ProblemInMemoryRepo) GetProblems() ([]domain.Problem, error) {
	return r.storage, nil
}

func (r *ProblemInMemoryRepo) GetProblemById(id int) (*domain.Problem, error) {
	for i := range r.storage {
		if r.storage[i].ID == id {
			return &r.storage[i], nil
		}
	}

	return nil, ErrNotFound
}

func (r *ProblemInMemoryRepo) GetProblemBySlug(slug string) (*domain.Problem, error) {
	for i := range r.storage {
		if r.storage[i].Slug == slug {
			return &r.storage[i], nil
		}
	}

	return nil, ErrNotFound
}
