package dto

type UserRequest struct {
	Name     string `json:"name" form:"name" binding:"required,alpha"`
	Username string `json:"username" form:"username" binding:"required,alphanum"`
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=8,max=25"`
}

type UserResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
