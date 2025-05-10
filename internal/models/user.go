package models

type User struct {
	ID           uint   `json:"id"       gorm:"primaryKey"`
	Name         string `json:"name"     binding:"required"`
	Email        string `json:"email"    gorm:"unique" binding:"required,email"`
	Age          int    `json:"age"      binding:"required,gte=0,lte=150"`
	PasswordHash string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type CreateUserRequest struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Age      int    `json:"age"      binding:"required,gte=0,lte=100"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"     binding:"required"`
	Email string `json:"email"    binding:"required,email"`
	Age   int    `json:"age"      binding:"required,gte=0,lte=150"`
}

type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
