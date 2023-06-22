package models

type Entity struct {
	Name         string  `json:"name"`
	EmailAddress *string `json:"email_address"`
	Username     *string `json:"username"`
	Password     string  `json:"password"`
	Icon         *string `json:"icon"`
	Description  *string `json:"description"`
	Type         int     `json:"type"`
}

const (
	github = iota
	facebook
	gmail
	linkedin
	instagram
	custom
)
