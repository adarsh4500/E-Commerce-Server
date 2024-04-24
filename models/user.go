package models

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"`
}

type Creds struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
