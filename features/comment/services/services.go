package services

import (
	"ALTA_BE_SOSMED/features/comment"
	"ALTA_BE_SOSMED/helper"
	"ALTA_BE_SOSMED/middlewares"
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
)

type service struct {
	m comment.ComModels
	v *validator.Validate
}

func NewComService(model comment.ComModels) comment.ComService {
	return &service{
		m: model,
		v: validator.New(),
	}
}

func (s *service) View(token *jwt.Token) ([]comment.Comment, error) {
	decodeHp := middlewares.DecodeToken(token)
	result, err := s.m.GetComByOwner(decodeHp)
	if err != nil {
		return []comment.Comment{}, err
	}

	return result, nil
}

func (s *service) AddCom(pemilik *jwt.Token, kegiatanBaru comment.Comment) (comment.Comment, error) {
	hp := middlewares.DecodeToken(pemilik)
	if hp == "" {
		log.Println("error decode token:", "token tidak ditemukan")
		return comment.Comment{}, errors.New("data tidak valid")
	}

	err := s.v.Struct(&kegiatanBaru)
	if err != nil {
		log.Println("error validasi", err.Error())
		return comment.Comment{}, err
	}

	result, err := s.m.InsertCom(hp, kegiatanBaru)
	if err != nil {
		return comment.Comment{}, errors.New(helper.ServerGeneralError)
	}

	return result, nil
}

func (s *service) Update(hp string, comID uint, newData comment.Comment) error {
	var updateValidate comment.Update
	updateValidate.Comment = newData.Comment
	updateValidate.Pemilik = newData.Pemiliks
	err := s.v.Struct(&updateValidate)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}

	_, err = s.m.Update(hp, comID, newData)
	if err != nil {
		return errors.New(helper.ServerGeneralError)
	}
	return nil
}

func (s *service) DeleteCom(deleteID comment.Comment) error {
	var deleteUser comment.DeleteCom
	deleteUser.ID = deleteID.ID
	err := s.v.Struct(&deleteUser)
	if err != nil {
		log.Println("error validasi", err.Error())
		return err
	}
	err = s.m.DeleteCom(deleteID)
	if err != nil {
		return errors.New(helper.ServerGeneralError)
	}
	return nil
}
