package service

type LoginService interface {
	Login(username, password string) bool
}

type loginService struct {
	authorizedUsername string
	authorizedPassword string
}

func NewLoginService() LoginService {
	return &loginService{
		authorizedUsername: "josorio",
		authorizedPassword: "pwd123",
	}
}

func (service *loginService) Login(username, password string) bool {
	return service.authorizedUsername == username && service.authorizedPassword == password
}
