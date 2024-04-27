package models

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Token struct {
	Uuid      string     `json:"uuid" gorm:"unique;type:uuid; column:uuid;default:uuid_generate_v4()"`
	UserUuid  string     `json:"user_uuid"`
	User      User       `gorm:"foreignKey:UserUuid;references:Uuid"`
	Token     string     `json:"token"`
	TokenType string     `json:"token_type"`
	ExpireAt  time.Time  `json:"expire_at"`
	IsUsed    *time.Time `json:"is_used"`
}

const (
	TOKEN_TYPE_FORGOT_PASSWORD = "forgot_password"
	TOKEN_TYPE_VERIFICATION    = "verification"
	TOKEN_TYPE_NEW_EMAIL       = "new_email"
	TOKEN_TYPE_DELETE_ACCOUNT  = "delete_accoun"
)

const DefaultTokenExpireTime = time.Minute * 30

func (t *Token) Validate() error {
	validate := validator.New()
	validate.RegisterStructValidationMapRules(TokenValidateRules, Token{})
	return validate.Struct(t)
}

var TokenValidateRules = map[string]string{
	"Token":     "required",
	"TokenType": "required,oneof=forgot_password verification new_email delete_account",
	"ExpireAt":  "required",
}
