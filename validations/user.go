package validations

import (
	"password-lock/models"
	"reflect"
)

func IsCompleteRegistrationRequestValid(personalQuestions []*models.UserPersonalQuestion) bool {

	for i, _ := range personalQuestions {

		rv := reflect.ValueOf(personalQuestions[i])
		rv = rv.Elem()

		for j := 0; j < rv.NumField(); j++ {

			if rv.Field(j).Kind() == reflect.Struct {
				continue
			}

			// If the field is a pointer, check if the dereferenced value is an empty string
			if rv.Field(j).Kind() == reflect.String {
				str := rv.Field(j).String()
				if str == "" {
					return false
				}
			}

		}
		return true
	}

	return true
}
