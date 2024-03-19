package handler

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password" `
}

type PictureRequest struct {
	Picture string `json:"picture" form:"picture"` 
}
