package dto

type (
	RegisterRequest struct {
		Name     string `json:"name" form:"name" binding:"required"`
		Username string `json:"username" form:"username" binding:"required,alphanum"`
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required,min=8,max=25"`
	}

	RegisterResponse struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	LoginRequest struct {
		Username string `json:"username" form:"username" binding:"required"`
		Password string `json:"password" form:"password" binding:"required,min=8,max=25"`
	}

	LoginResponse struct {
		ExpiredIn    int64  `json:"expired_in"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
)
