package models

import (
	"encoding/json"
	"password-lock/db"
)

type PersonalQuestion struct {
	Uuid     string `json:"uuid" gorm:"unique;type:uuid; column:uuid;default:uuid_generate_v4()"`
	Question string `json:"question"`
}

type UserPersonalQuestion struct {
	PersonalQuestion     PersonalQuestion `gorm:"foreignKey:PersonalQuestionUuid;references:Uuid;"`
	PersonalQuestionUuid string           `json:"personal_question_uuid"`
	Answer               string           `json:"answer"`
	UserUuid             string           `json:"user_uuid"`
}

func (p PersonalQuestion) TableName() string {
	return db.PERSONAL_QUESTIONS_TABLE
}

func (p UserPersonalQuestion) TableName() string {
	return db.USER_PERSONAL_QUESTIONS_TABLE
}

func (p UserPersonalQuestion) MarshalJSON() ([]byte, error) {
	var tmp struct {
		PersonalQuestion PersonalQuestion `json:"personal_question"`
	}

	tmp.PersonalQuestion = p.PersonalQuestion

	return json.Marshal(&tmp)
}
