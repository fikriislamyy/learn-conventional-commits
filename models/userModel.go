package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct{
	gorm.Model
	ID 			uuid.UUID `gorm:"type:uuid;primary_key"`
	Firstname 	string
	Lastname 	string
	Email 		string `gorm:"unique"`
	Password 	string
	Phone 		string `gorm:"unique"`
	Role 		string `gorm:"default:user"`
	CreatedAt 	time.Time
	UpdatedAt 	time.Time
}