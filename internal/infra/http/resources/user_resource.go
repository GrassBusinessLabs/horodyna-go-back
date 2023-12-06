package resources

import "boilerplate/internal/domain"

type UserDto struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UsersDto struct {
	Items []UserDto `json:"items"`
	Total uint64    `json:"total"`
	Pages uint      `json:"pages"`
}

type AuthDto struct {
	Token string  `json:"token"`
	User  UserDto `json:"user"`
}

func (d UserDto) DomainToDto(user domain.User) UserDto {
	return UserDto{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}

func (d UserDto) DomainToDtoCollection(users domain.Users) UsersDto {
	result := make([]UserDto, len(users.Items))

	for i := range users.Items {
		result[i] = d.DomainToDto(users.Items[i])
	}

	return UsersDto{Items: result, Pages: users.Pages, Total: users.Total}
}

func (d AuthDto) DomainToDto(token string, user domain.User) AuthDto {
	var userDto UserDto
	return AuthDto{
		Token: token,
		User:  userDto.DomainToDto(user),
	}
}
