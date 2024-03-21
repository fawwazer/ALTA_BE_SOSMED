package handler

type LoginResponse struct {
	Email string `json:"email"`
	Nama  string `json:"nama"`
	Token string `json:"token"`
}

type ProfileResponse struct {
	UserID    int
	Nama      string
	Email     string
	Picture   string
	Tgl_lahir string
	Gender    bool
}
