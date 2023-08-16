package models

type Entity struct {
	Name         string  `json:"name"`
	EmailAddress *string `json:"email_address"`
	Username     *string `json:"username"`
	Password     string  `json:"password"`
	IconPath     *string `json:"icon_path"`
	Description  *string `json:"description"`
	Type         int     `json:"type"`
	UserUuid     string  `json:"user_uuid"`
	SecretKey    string  `json:"secret_key" gorm:"-"`
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
