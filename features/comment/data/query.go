package data

import (
	"ALTA_BE_SOSMED/features/comment"
	"errors"

	"gorm.io/gorm"
)

type ComModel struct {
	connection *gorm.DB
}

func New(db *gorm.DB) comment.ComModels {
	return &ComModel{
		connection: db,
	}
}

func (tm *ComModel) InsertCom(pemilik string, inputCom comment.Comment) (comment.Comment, error) {
	var inputProcess = comment.Comment{
		Comment:  inputCom.Comment,
		Pemiliks: pemilik,
	}
	if err := tm.connection.Create(&inputProcess).Error; err != nil {
		return comment.Comment{}, err
	}

	return inputProcess, nil
}

func (tm *ComModel) Update(pemilik string, comID uint, data comment.Comment) (comment.Comment, error) {
	var qry = tm.connection.Where("pemiliks = ? AND id = ?", pemilik, comID).Updates(data)
	if err := qry.Error; err != nil {
		return comment.Comment{}, err
	}

	if qry.RowsAffected == 0 {
		return comment.Comment{}, errors.New("tidak ada data yg di update")
	}

	return data, nil
}

func (tm *ComModel) GetComByOwner(pemilik string) ([]comment.Comment, error) {
	var result []comment.Comment
	if err := tm.connection.Where("pemiliks = ?  ", pemilik).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (m *ComModel) DeleteCom(deleteID comment.Comment) error {
	var data = m.connection.Delete(&Comment{}, deleteID)
	if err := data.Error; err != nil {
		return err
	}
	if data.RowsAffected == 0 {
		return errors.New("tidak ada data yg dihapus")
	}
	return nil
}
