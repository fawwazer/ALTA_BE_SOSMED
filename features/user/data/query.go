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

func (m *UserModel) Login(email string) (user.User, error) {
	var result user.User
	if err := m.Connection.Where("email = ? ", email).First(&result).Error; err != nil {
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

func (um *UserModel) GetLastUserID() (int, error) {
	var lastUser User

	// Query untuk mendapatkan userID terakhir berdasarkan id terbesar
	if err := um.Connection.Order("user_id desc").First(&lastUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Tabel kosong, return 0 sebagai userID pertama
			return 0, nil
		}
		return 0, err
	}

	return lastUser.UserID, nil
}

func (um *UserModel) UploadPictureURL(userID int, picture string) error {
	if err := um.Connection.Model(&User{}).Where("user_id = ?", userID).Update("picture", picture).Error; err != nil {
		return err
	}
	return nil
}
