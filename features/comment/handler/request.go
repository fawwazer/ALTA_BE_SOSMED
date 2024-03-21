package handler

type ComRequest struct {
	ID       string `json:"id"`
	Comment  string `json:"comment" form:"comment" validate:"required"`
	Pemiliks string `gorm:"type:varchar(30)"`
}
