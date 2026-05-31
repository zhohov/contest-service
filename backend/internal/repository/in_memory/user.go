package repository

import (
	"github.com/zhohov/contest-service/backend/internal/domain"
)

type UserInMemoryRepo struct {
	id      int
	storage []domain.User
}

func NewUserInMemoryRepo(storage []domain.User) *UserInMemoryRepo {
	return &UserInMemoryRepo{
		id:      len(storage) + 1,
		storage: storage,
	}
}

func (r *UserInMemoryRepo) CreateUser(user *domain.User) (*domain.User, error) {
	user.ID = r.id
	r.id++

	r.storage = append(r.storage, *user)

	return user, nil
}

func (r *UserInMemoryRepo) GetUsers() ([]domain.User, error) {
	return r.storage, nil
}

func (r *UserInMemoryRepo) GetUserById(id int) (*domain.User, error) {
	for i := range r.storage {
		if r.storage[i].ID == id {
			return &r.storage[i], nil
		}
	}

	return nil, ErrNotFound
}

func (r *UserInMemoryRepo) GetUserByEmail(email string) (*domain.User, error) {
	for i := range r.storage {
		if r.storage[i].Email == email {
			return &r.storage[i], nil
		}
	}

	return nil, ErrNotFound
}
