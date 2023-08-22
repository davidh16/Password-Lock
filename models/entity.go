package models

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
	SecretKey    string  `json:"secret_key,omitempty" gorm:"-"`
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
	return "entities"
}

func (e *Entity) Merge(data *Entity) {

	if data.Name != "" {
		e.Name = data.Name
	}

	if data.EmailAddress != nil {
		e.EmailAddress = data.EmailAddress
	}

	if data.Username != nil {
		e.Username = data.Username
	}

	if data.IconPath != "" {
		e.IconPath = data.IconPath
	}

	if data.Description != nil {
		e.Description = data.Description
	}

	if data.Type != 0 {
		e.Type = data.Type
	}

	if data.Password != "" {
		e.Password = data.Password
	}
}
