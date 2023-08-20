package models

type Entity struct {
	Name         string  `json:"name,omitempty"`
	EmailAddress *string `json:"email_address,omitempty"`
	Username     *string `json:"username,omitempty"`
	Password     string  `json:"password,omitempty"`
	IconPath     *string `json:"icon_path,omitempty"`
	Description  *string `json:"description,omitempty"`
	Type         int     `json:"type,omitempty"`
	UserUuid     string  `json:"user_uuid,omitempty"`
	SecretKey    string  `json:"secret_key,omitempty" gorm:"-"`
}

const (
	github = iota
	facebook
	gmail
	linkedin
	instagram
	custom
)

func (e Entity) TableName() string {
	return "entities"
}
