package controllers

import (
	"boilerplate/internal/app"
	"boilerplate/internal/domain"
	"boilerplate/internal/infra/http/requests"
	"boilerplate/internal/infra/http/resources"
	"errors"
	"log"
	"net/http"
)

type UserController struct {
	userService app.UserService
}

func NewUserController(us app.UserService) UserController {
	return UserController{
		userService: us,
	}
}

func (c UserController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
		}

		user, err = c.userService.Save(user)
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}

		var userDto resources.UserDto
		Created(w, userDto.DomainToDto(user))
	}
}

func (c UserController) FindMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)
		Success(w, resources.UserDto{}.DomainToDto(user))
	}
}

func (c UserController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.UpdateUserRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}

		u := r.Context().Value(UserKey).(domain.User)
		user, err = c.userService.Update(u)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		var userDto resources.UserDto
		Success(w, userDto.DomainToDto(user))
	}
}

func (c UserController) SetPhoneNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userPhoneNumber, err := requests.Bind(r, requests.SetPhoneNumberRequest{}, domain.User{})
		if err != nil {
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}
		user := r.Context().Value(UserKey).(domain.User)
		if user.PhoneNumber != nil {
			err = errors.New("user already have a phone number")
			log.Printf("UserController: %s", err)
			BadRequest(w, err)
			return
		}
		user.PhoneNumber = userPhoneNumber.PhoneNumber
		user, err = c.userService.Update(user)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}
		var userDto resources.UserDto
		Success(w, userDto.DomainToDto(user))
	}
}

func (c UserController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(UserKey).(domain.User)

		err := c.userService.Delete(u.Id)
		if err != nil {
			log.Printf("UserController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
