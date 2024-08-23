package models

import (
	"password-lock/db"
	"reflect"
)

type Entity struct {
	Uuid         string  `json:"uuid" gorm:"unique;type:uuid; column:uuid;default:uuid_generate_v4()"`
	Name         string  `json:"name,omitempty"`
	EmailAddress *string `json:"email_address,omitempty"`
	Username     *string `json:"username,omitempty"`
	Password     string  `json:"password,omitempty"`
	IconPath     string  `json:"icon_path,omitempty"`
	Description  *string `json:"description,omitempty"`
	Type         int     `json:"type,omitempty"`
	UserUuid     string  `json:"user_uuid,omitempty"`
	User         User    `gorm:"constraint:OnDelete:CASCADE;foreignKey:Uuid;references:UserUuid" json:"-"`
}

var TypeMap = map[int]string{
	1: "github",
	2: "facebook",
	3: "gmail",
	4: "linkedin",
	5: "instagram",
	6: "custom",
}

func (e *Entity) HidePassword() *Entity {
	e.Password = ""
	return e
}

func (e *Entity) TableName() string {
	return db.ENTITIES_TABLE
}

func (e *Entity) Merge(data *Entity) {

	if !reflect.DeepEqual(e.Name, data.Name) {
		e.Name = data.Name
	}

	if !reflect.DeepEqual(e.EmailAddress, data.EmailAddress) {
		e.EmailAddress = data.EmailAddress
	}

	if !reflect.DeepEqual(e.Username, data.Username) {
		e.Username = data.Username
	}

	if !reflect.DeepEqual(e.IconPath, data.IconPath) {
		e.IconPath = data.IconPath
	}

	if !reflect.DeepEqual(e.Description, data.Description) {
		e.Description = data.Description
	}

	if !reflect.DeepEqual(e.Type, data.Type) {
		e.Type = data.Type
	}

	if !reflect.DeepEqual(e.IconPath, data.IconPath) {
		e.IconPath = data.IconPath
	}

}
