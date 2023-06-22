package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Uuid         string    `json:"uuid"`
	EmailAddress string    `json:"email_address"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
}

func (u *User) Validate() error {
	v := validator.New()
	v.RegisterStructValidationMapRules(ValidationRules, User{})
	err := v.Struct(u)
	if err != nil {
		return err
	}
	return nil
}

var ValidationRules = map[string]string{
	"EmailAddress": "required,email",
	"Password":     "required,min=8",
}

// BeforeCreate - Gorm hook that encrypts password before saving user to database
func (u *User) BeforeCreate(tx *gorm.DB) error {

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 20)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}
