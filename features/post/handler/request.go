package handler

type PostRequest struct {
	Posting string `json:"posting" form:"posting"`
	Picture string `json:"picture" form:"picture"`
}
