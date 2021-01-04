package use_cases

import (
	"taurus/domains/models"
	"taurus/domains/repositories"
)

// 外側から内側を操作するときの抽象
type UserUsecase interface {
	Store(name, email, password string) (*models.User, error)
	List() ([]*models.User, error)
	Find(id string) (*models.User, error)
	Update(id, name, email, password string) (*models.User, error)
	Delete(id string) error
}

// 実態はこっち
type userUsecase struct {
	repo repositories.UserRepository
}

// returnの型を抽象として実態を返す
func NewUserUsecase(repo repositories.UserRepository) UserUsecase {
	return &userUsecase{repo}
}

func (u *userUsecase) Store(name, email, password string) (*models.User, error) {
	user, err := u.repo.Store(name, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) List() ([]*models.User, error) {
	users, err := u.repo.List()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) Find(id string) (*models.User, error) {
	user, err := u.repo.Find(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Update(id, name, email, password string) (*models.User, error) {
	user, err := u.repo.Update(id, name, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Delete(id string) error {
	if err := u.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
