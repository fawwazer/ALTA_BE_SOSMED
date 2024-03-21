package handler

type PostRequest struct {
	Posting string `json:"kegiatan" form:"kegiatan"`
}
