package dto

type RegisterInput struct {
	Email    string `json:"email" form:"email" binding:"required,email,max=255"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=255"`
}

type LoginInput struct {
	Email    string `json:"email" form:"email" binding:"required,email,max=255"`
	Password string `json:"password" form:"password" binding:"required,max=255"`
}
