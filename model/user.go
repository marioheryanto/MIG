package model

type User struct {
	Id       uint   `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"password"`
}

type UserRequest struct {
	Id       uint   `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
