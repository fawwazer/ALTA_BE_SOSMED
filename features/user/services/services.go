package services

import (
	"ALTA_BE_SOSMED/features/user"

	"ALTA_BE_SOSMED/helper"
	"ALTA_BE_SOSMED/middlewares"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	// "github.com/labstack/gommon/email"
)

type service struct {
	model user.UserModel
	pm    helper.PasswordManager
	v     *validator.Validate
}

func NewService(m user.UserModel) user.UserService {
	return &service{
		model: m,
		pm:    helper.NewPasswordManager(),
		v:     validator.New(),
	}
}

func (s *service) Register(newData user.User) error {
	var registerValidate user.Register

	// Mencari UserID terakhir dari database
	lastUserID, error := s.model.GetLastUserID()
	if error != nil {
		return errors.New(helper.ServiceGeneralError)
	}

	// Menentukan UserID untuk pengguna baru
	newUserID := lastUserID + 1
	newData.UserID = newUserID

	registerValidate.UserID = newData.UserID
	registerValidate.Email = newData.Email
	registerValidate.Nama = newData.Nama
	registerValidate.Password = newData.Password
	registerValidate.Tgl_lahir = newData.Tgl_lahir
	registerValidate.Gender = newData.Gender
	err := s.v.Struct(&registerValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	newPassword, err := s.pm.HashPassword(newData.Password)
	if err != nil {
		return errors.New(helper.ServiceGeneralError)
	}
	newData.Password = newPassword

	err = s.model.AddUser(newData)
	if err != nil {
		return errors.New(helper.ServerGeneralError)
	}
	return nil
}

func (s *service) Login(loginData user.User) (user.User, string, error) {
	var loginValidate user.Login
	loginValidate.Email = loginData.Email
	loginValidate.Password = loginData.Password
	err := s.v.Struct(&loginValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return user.User{}, "", err
	}

	dbData, err := s.model.Login(loginValidate.Email)
	if err != nil {
		log.Println("error login model", err.Error())
		return user.User{}, "", errors.New(helper.UserCredentialError) //
	}

	err = s.pm.ComparePassword(loginValidate.Password, dbData.Password)
	if err != nil {
		log.Println("error compare", err.Error())
		return user.User{}, "", errors.New(helper.UserCredentialError)
	}

	token, err := middlewares.GenerateJWT(dbData.Email)
	if err != nil {
		log.Println("error generate", err.Error())
		return user.User{}, "", errors.New(helper.ServiceGeneralError)
	}

	return dbData, token, nil
}

func (s *service) Profile(token *jwt.Token, userID uint) (user.User, error) {
	decodeEmail := middlewares.DecodeToken(token)
	if decodeEmail == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return user.User{}, errors.New("data tidak valid")
	}

	result, err := s.model.GetUserByID(userID)
	if err != nil {
		return user.User{}, err
	}

	if result.Email != decodeEmail {
		return user.User{}, errors.New("anda tidak diizinkan mengakses profil pengguna lainn")
	}

	return result, nil
}

// saveUploadedFile function to handle file uploads.
func (s *service) SaveUploadedFile(file *multipart.FileHeader, path string) error {
	// Open the uploaded file.
	src, err := file.Open()
	if err != nil {
		log.Print("file open error :", err.Error())
		return err
	}
	defer src.Close()

	// Create a destination file for the uploaded content.
	dst, err := os.Create(path)
	if err != nil {
		log.Print("file create error :", err.Error())
		return err
	}
	defer dst.Close()

	// Copy the uploaded content to the destination file.
	if _, err = io.Copy(dst, src); err != nil {
		log.Print("file copy error :", err.Error())
		return err
	}

	return nil
}

func (s *service) UpdateProfile(userID int, token *jwt.Token, newData user.User) error {
	email := middlewares.DecodeToken(token)
	if email == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return errors.New("data tidak valid")
	}

	user, error := s.model.GetUserByID(uint(userID))
	if error != nil {
		log.Println("error getting user:", error.Error())
		return error
	}

	if user.Email != email {
		log.Println("error get account:", "user tidak sesuai")
		return errors.New("user tidak sesuai")
	}

	err := s.v.Struct(&newData)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	// membuat map untuk menampung kolom yang akan diperbarui bersama dengan nilainya
	updateFields := make(map[string]interface{})

	// tentukan kolom yang ingin diperbarui dan tambahkan ke dalam map
	if newData.Nama != "" {
		updateFields["nama"] = newData.Nama
	}
	if newData.Email != "" {
		updateFields["email"] = newData.Email
	}
	if newData.Tgl_lahir != "" {
		updateFields["tgl_lahir"] = newData.Tgl_lahir
	}
	if newData.Picture != "" {
		updateFields["picture"] = newData.Picture
	}
	updateFields["gender"] = newData.Gender

	err = s.model.Update(userID, updateFields, email)
	if err != nil {
		log.Print("error update to model: ", err.Error())
		return errors.New(helper.ServerGeneralError)
	}

	return nil
}

func (s *service) GetPicture(token *jwt.Token) (user.User, error) {
	email := middlewares.DecodeToken(token)
	if email == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return user.User{}, errors.New("data tidak valid")
	}

	result, err := s.model.GetPicture(email)
	if err != nil {
		log.Println("error getting user:", err.Error())
		return user.User{}, err
	}

	return result, err
}

func (s *service) DeleteAccount(userID uint, token *jwt.Token) error {
	email := middlewares.DecodeToken(token)
	if email == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return errors.New("data tidak valid")
	}

	user, err := s.model.GetUserByID(userID)
	if err != nil {
		log.Println("error getting user:", err.Error())
		return err
	}

	if user.Email != email {
		log.Println("error deleting account:", "user tidak sesuai")
		return errors.New("user tidak sesuai")
	}

	error := s.model.Delete(userID, email)
	if error != nil {
		log.Print("error delete to model: ", error.Error())
		return errors.New(helper.ServerGeneralError)
	}

	return nil
}
