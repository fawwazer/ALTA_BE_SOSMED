package data

import (
	"ALTA_BE_SOSMED/features/user"

	"gorm.io/gorm"
)

type UserModel struct {
	Connection *gorm.DB
}

func New(db *gorm.DB) user.UserModel {
	return &UserModel{
		Connection: db,
	}
}

func (um *UserModel) AddUser(newData user.User) error {
	err := um.Connection.Create(&newData).Error
	if err != nil {
		return err
	}
	return nil
}

func (um *UserModel) CekUser(email string) bool {
	var data User
	if err := um.Connection.Where("Email = ?", email).First(&data).Error; err != nil {
		return false
	}
	return true
}

func (um *UserModel) UpdateUser(email string, data user.User) error {
	if err := um.Connection.Model(&data).Where("Email = ?", email).Update("Name", data.Nama).Update("Password", data.Password).Error; err != nil {
		return err
	}
	return nil
}

func (um *UserModel) Login(email string) (user.User, error) {
	var result user.User
	if err := um.Connection.Where("Email = ? ", email).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

func (um *UserModel) GetUserByEmail(email string) (user.User, error) {
	var result user.User
	if err := um.Connection.Where("Email = ?", email).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}
