package requests

import (
	"boilerplate/internal/domain"
)

type RegisterRequest struct {
	Name        string `json:"name" validate:"required,gte=1,max=40"`
	Email       string `json:"email" validate:"email"`
	Password    string `json:"password" validate:"required,alphanum,gte=4,max=20"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type AuthRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required,alphanum,gte=4"`
}

type EmailAuthRequest struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required,alphanum,gte=4"`
}

type UpdateUserRequest struct {
	Name    string `json:"name" validate:"required,gte=1,max=40"`
	From    string `json:"from" validate:"required,email"`
	To      string `json:"to" validate:"required,email"`
	GoTrans bool   `json:"go_trans" validate:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required,alphanum,gte=4"`
	NewPassword string `json:"newPassword" validate:"required,alphanum,gte=4"`
}

type SetPhoneNumberRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

func (r UpdateUserRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Name: r.Name,
	}, nil
}

func (r RegisterRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Email:       r.Email,
		Password:    r.Password,
		PhoneNumber: &r.PhoneNumber,
		Name:        r.Name,
	}, nil
}

func (r AuthRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		PhoneNumber: &r.PhoneNumber,
		Password:    r.Password,
	}, nil
}

func (r EmailAuthRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Email:    r.Email,
		Password: r.Password,
	}, nil
}

func (r ChangePasswordRequest) ToDomainModel() (interface{}, error) {
	return domain.ChangePassword{
		OldPassword: r.OldPassword,
		NewPassword: r.NewPassword,
	}, nil
}

func (r SetPhoneNumberRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		PhoneNumber: &r.PhoneNumber,
	}, nil
}
