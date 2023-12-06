package requests

import (
	"boilerplate/internal/domain"
)

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,gte=1,max=40"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,gte=4,max=20"`
}

type AuthRequest struct {
	Email    string `json:"email"  validate:"required,email"`
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

func (r UpdateUserRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Name: r.Name,
	}, nil
}

func (r RegisterRequest) ToDomainModel() (interface{}, error) {
	return domain.User{
		Email:    r.Email,
		Password: r.Password,
		Name:     r.Name,
	}, nil
}

func (r AuthRequest) ToDomainModel() (interface{}, error) {
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
