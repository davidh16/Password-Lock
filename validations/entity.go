package validations

import (
	"errors"
	"password-lock/models"
)

func ValidateCreateEntityRequest(entity models.Entity) error {
	if entity.EmailAddress != nil && *entity.EmailAddress == "" {
		return errors.New("validation error : invalid email address")
	}

	if entity.Username != nil && *entity.Username == "" {
		return errors.New("validation error : invalid username")
	}

	if entity.Type > 6 && entity.Type < 1 {
		return errors.New("validation error : invalid type")
	}

	if entity.Type == 6 && entity.Name == "" {
		return errors.New("validation error : invalid name")
	}

	if entity.Description != nil {
		if *entity.Description == "" {
			return errors.New("validation error : invalid description")
		}
	}

	return nil
}
