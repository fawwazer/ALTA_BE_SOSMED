package services

import (
	"ALTA_BE_SOSMED/features/user"
	"ALTA_BE_SOSMED/helper"
	"ALTA_BE_SOSMED/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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
	registerValidate.Email = newData.Email
	registerValidate.Nama = newData.Nama
	registerValidate.Password = newData.Password
	registerValidate.Tgl_lahir = newData.Tgl_lahir
	registerValidate.Gender = newData.Gender
	registerValidate.Alamat = newData.Alamat
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
		return user.User{}, "", err
	}

	err = s.pm.ComparePassword(loginValidate.Password, dbData.Password)
	if err != nil {
		return user.User{}, "", errors.New(helper.UserCredentialError)
	}

	token, err := middlewares.GenerateJWT(dbData.Email)
	if err != nil {
		return user.User{}, "", errors.New(helper.ServiceGeneralError)
	}

	return dbData, token, nil
}

func (s *service) Profile(token *jwt.Token) (user.User, error) {
	decodeEmail := middlewares.DecodeToken(token)
	result, err := s.model.GetUserByEmail(decodeEmail)
	if err != nil {
		return user.User{}, err
	}
	return result, nil
}
