package repositories

import (
	"database/sql"
	"taurus/domains/models"
	"taurus/domains/repositories"
	"taurus/infrastructures"
	"time"

	uuid "github.com/satori/go.uuid"
)

// domains/repositories の実装
type UserRepository struct {
	infrastructures.Database
}

func NewUserRepository(database infrastructures.Database) repositories.UserRepository {
	return &UserRepository{database}
}

func (r *UserRepository) Store(name, email, password string) (*models.User, error) {
	var newUser = &models.User{
		Name:  name,
		Email: email,
		// 気が向いたらbcrypt導入する
		Password: password,
	}
	if err := r.Database.Create(newUser); err != nil {
		return nil, err
	}
	return newUser, nil
}

func (r *UserRepository) List() ([]*models.User, error) {
	var users []*models.User
	targetModel := &models.User{IsDeleted: true}
	if err := r.Database.Not(targetModel).Find(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) Find(id string) (*models.User, error) {
	var user = &models.User{}
	targetModel := &models.User{ID: uuid.FromStringOrNil(id)}
	if err := r.Database.Where(targetModel).First(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(id, name, email, password string) (*models.User, error) {
	var user = &models.User{
		Name:      name,
		Email:     email,
		Password:  password,
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}
	targetModel := &models.User{ID: uuid.FromStringOrNil(id)}
	if err := r.Database.Where(targetModel).Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Delete(id string) error {
	user := &models.User{IsDeleted: true}
	targetModel := &models.User{ID: uuid.FromStringOrNil(id)}
	if err := r.Database.Where(targetModel).Update(user); err != nil {
		return err
	}
	return nil
}
