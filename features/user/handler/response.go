package handler

type LoginResponse struct {
	Email string `json:"email"`
	Nama  string `json:"nama"`
	Token string `json:"token"`
}
