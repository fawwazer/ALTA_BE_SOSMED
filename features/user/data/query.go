package data

import "gorm.io/gorm"

type UserModel struct {
	Connection *gorm.DB
}

func (um *UserModel) AddUser(newData User) error {
	err := um.Connection.Create(&newData).Error
	if err != nil {
		return err
	}
	return nil
}

func New(db *gorm.DB) *UserModel {
	return &UserModel{
		Connection: db,
	}
}

func (um *UserModel) CekUser(email string) bool {
	var data User
	if err := um.Connection.Where("Email = ?", email).First(&data).Error; err != nil {
		return false
	}
	return true
}

func (um *UserModel) Update(email string, data User) error {
	if err := um.Connection.Model(&data).Where("Email = ?", email).Update("Name", data.Nama).Update("Password", data.Password).Error; err != nil {
		return err
	}
	return nil
}

func (um *UserModel) Login(email string, password string) (User, error) {
	var result User
	if err := um.Connection.Where("Email = ? AND Password = ?", email, password).First(&result).Error; err != nil {
		return User{}, err
	}
	return result, nil
}

func (um *UserModel) GetUserByEmail(email string) (User, error) {
	var result User
	if err := um.Connection.Where("Email = ?", email).First(&result).Error; err != nil {
		return User{}, err
	}
	return result, nil
}
