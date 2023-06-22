package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
	"Username": "required",
	"Password": "required,min=8",
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
