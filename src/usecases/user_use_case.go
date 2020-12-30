package usecases

import (
	"taurus/domains/models"
	"taurus/domains/repositories"
)

// 外側から内側を操作するときの抽象
type UserUsecase interface {
	PostUser(name, email, password string) (*models.User, error)
	ListUser() ([]*models.User, error)
	FindUser(id string) (*models.User, error)
	UpdateUser(id, name, email, password string) (*models.User, error)
	DeleteUser(id string) error
}

// 実態はこっち
type userUsecase struct {
	repo repositories.UserRepository
}

// returnの型を抽象として実態を返す
func NewUserUsecase(repo repositories.UserRepository) UserUsecase {
	return &userUsecase{repo}
}

func (u *userUsecase) PostUser(name, email, password string) (*models.User, error) {
	user, err := u.repo.Store(name, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) ListUser() ([]*models.User, error) {
	users, err := u.repo.Scan()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userUsecase) FindUser(id string) (*models.User, error) {
	user, err := u.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) UpdateUser(id, name, email, password string) (*models.User, error) {
	user, err := u.repo.Update(id, name, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) DeleteUser(id string) error {
	if err := u.repo.Delete(id); err != nil {
		return err
	}
	return nil
}
