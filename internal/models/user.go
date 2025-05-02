package models

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" gorm:"unique" binding:"required,email"`
	Age          int    `json:"age" binding:"required,gte=0,lte=150"`
	PasswordHash string `json:"password" binding:"required"` // пароль приходит в теле, но не возвращается
}
