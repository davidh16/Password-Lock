package validations

import (
	"errors"
	"github.com/samber/lo"
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

func ValidatePersonalQuestionsAnswers(userPersonalQuestions []models.UserPersonalQuestion, personalQuestionsFromRequest []models.UserPersonalQuestion) error {
	for _, personalQuestion := range userPersonalQuestions {
		var answerFromRequest string
		_, exist := lo.Find(personalQuestionsFromRequest, func(item models.UserPersonalQuestion) bool {
			if item.PersonalQuestionUuid == personalQuestion.PersonalQuestionUuid {
				answerFromRequest = item.Answer
				return true
			}
			return false
		})

		if exist {
			if personalQuestion.Answer != answerFromRequest {
				return errors.New("we could not verify your identity")
			}
		} else {
			return errors.New("we could not verify your identity")
		}
	}

	return nil
}
