package models

type Secret struct {
	Host     string `json:"host"`
	UserName string `json:"username"`
	Password string `json:"password"`
	JwtSign  string `json:"jwtsign"`
	Database string `json:"database"`
}
