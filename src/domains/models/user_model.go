package models

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

// User is m_users table model
type User struct {
	ID        uuid.UUID `gorm:"primary_key;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:text;not null"`
	Email     string    `gorm:"type:text;not null"`
	Password  string    `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	IsDeleted bool `gorm:"default:false;not null"`
}

// structの名前とtable名が違うので置換させるため
func (User) TableName() string {
	return "m_users"
}
