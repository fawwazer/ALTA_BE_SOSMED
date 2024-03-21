package data

import (
	"ALTA_BE_SOSMED/features/user"
	"errors"
	"log"

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

func (um *UserModel) GetUserByID(userID uint) (user.User, error) {
	var result user.User
	if err := um.Connection.Where("user_id = ?", userID).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

func (um *UserModel) GetLastUserID() (int, error) {
	var lastUser User

	// query untuk mendapatkan userID terakhir berdasarkan id terbesar
	if err := um.Connection.Order("user_id desc").First(&lastUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// tabel kosong, return 0 sebagai userID pertama
			return 0, nil
		}
		return 0, err
	}

	return lastUser.UserID, nil
}

func (um *UserModel) GetPicture(email string) (user.User, error) {
	var result user.User
	if err := um.Connection.Where("email = ?", email).First(&result).Error; err != nil {
		return user.User{}, err
	}
	return result, nil
}

func (um *UserModel) Update(userID int, updateFields map[string]interface{}, email string) error {
	var query = um.Connection.Model(&User{}).Where("user_id = ? AND email = ?", userID, email).Updates(updateFields)
	if err := query.Error; err != nil {
		log.Print("error to database :", err.Error())
		return err
	}
	if query.RowsAffected < 1 {
		return errors.New("no data affected")
	}
	return nil
}

func (um *UserModel) Delete(userID uint, email string) error {
	// if err := um.Connection.Unscoped().Where("user_id = ?", userID).Delete(userID).Error; err != nil {
	// 	log.Print("error to database :", err.Error())
	// 	return err
	// }
	if err := um.Connection.Model(&User{}).Where("user_id = ? AND email = ?", userID, email).Delete(userID).Error; err != nil {
		log.Print("error to database :", err.Error())
		return err
	}
	return nil
}
