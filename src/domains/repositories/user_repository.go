package repositories

import "taurus/domains/models"

// UserRepository データ操作の抽象
type UserRepository interface {
	Store(name, email, password string) (*models.User, error)
	List() ([]*models.User, error)
	Find(id string) (*models.User, error)
	Update(id, name, email, password string) (*models.User, error)
	Delete(id string) error
}
