package models

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Uuid              string                 `json:"uuid" gorm:"unique;type:uuid; column:uuid;default:uuid_generate_v4()"`
	EmailAddress      string                 `json:"email_address"`
	Password          string                 `json:"password"`
	CreatedAt         time.Time              `json:"created_at"`
	Active            bool                   `json:"active"`
	Completed         bool                   `json:"completed"`
	PersonalQuestions []UserPersonalQuestion `gorm:"foreignKey:UserUuid;references:Uuid"`
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

// BeforeSave - Gorm hook that encrypts password before saving user to database
func (u *User) BeforeSave(tx *gorm.DB) error {

	if u.Password != "" {
		var encryptPassword bool
		value, ok := tx.Get("encrypt-password")
		if ok {
			encryptPassword = value.(bool)
		}
		if encryptPassword {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
			if err != nil {
				return err
			}

			u.Password = string(hashedPassword)
		}
	}

	return nil
}

func (u *User) TableName() string {
	return "users"
}
