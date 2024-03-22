package data

import (
	"ALTA_BE_SOSMED/features/post"
	"errors"

	"gorm.io/gorm"
)

type model struct {
	connection *gorm.DB
}

func New(db *gorm.DB) post.PostModel {
	return &model{
		connection: db,
	}
}

func (tm *model) InsertPost(pemilik string, postingBaru post.Post) (post.Post, error) {
	var inputProcess = Post{Posting: postingBaru.Posting, Pemilik: pemilik, Picture: postingBaru.Picture}
	if err := tm.connection.Create(&inputProcess).Error; err != nil {
		return post.Post{}, err
	}
	return post.Post{Posting: inputProcess.Posting, Picture: inputProcess.Picture}, nil
}

func (tm *model) UpdatePost(pemilik string, postID uint, data post.Post) (post.Post, error) {
	var qry = tm.connection.Where("pemilik = ? AND id = ?", pemilik, postID).Updates(data)
	if err := qry.Error; err != nil {
		return post.Post{}, err
	}

	if qry.RowsAffected < 1 {
		return post.Post{}, errors.New("no data affected")
	}

	return data, nil
}

func (tm *model) GetPostByOwner(pemilik string) ([]post.Post, error) {
	var result []post.Post
	if err := tm.connection.Where("pemilik = ?", pemilik).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}
