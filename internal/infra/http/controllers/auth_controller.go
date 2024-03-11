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

type AuthController struct {
	authService app.AuthService
	userService app.UserService
}

func NewAuthController(as app.AuthService, us app.UserService) AuthController {
	return AuthController{
		authService: as,
		userService: us,
	}
}

func (c AuthController) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.RegisterRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, errors.New("invalid request body"))
			return
		}

		user, token, err := c.authService.Register(user)
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		var authDto resources.AuthDto
		Success(w, authDto.DomainToDto(token, user))
	}
}

func (c AuthController) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.AuthRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		u, token, err := c.authService.Login(user)
		if err != nil {
			Unauthorized(w, err)
			return
		}

		var authDto resources.AuthDto
		Success(w, authDto.DomainToDto(token, u))
	}
}

func (c AuthController) LoginWithEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := requests.Bind(r, requests.EmailAuthRequest{}, domain.User{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}

		u, token, err := c.authService.LoginWithEmail(user)
		if err != nil {
			Unauthorized(w, err)
			return
		}

		var authDto resources.AuthDto
		Success(w, authDto.DomainToDto(token, u))
	}
}

func (c AuthController) Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := r.Context().Value(SessKey).(domain.Session)
		err := c.authService.Logout(sess)
		if err != nil {
			log.Print(err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

func (c AuthController) ChangePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := requests.Bind(r, requests.ChangePasswordRequest{}, domain.ChangePassword{})
		if err != nil {
			log.Printf("AuthController: %s", err)
			BadRequest(w, err)
			return
		}
		sess := r.Context().Value(SessKey).(domain.Session)
		user := r.Context().Value(UserKey).(domain.User)

		err = c.authService.ChangePassword(user, req, sess)
		if err != nil {
			log.Printf("AuthController: %s", err)
			InternalServerError(w, err)
			return
		}

		Ok(w)
	}
}
