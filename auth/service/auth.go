package service

type AuthService struct{}

func NewAuthService() (*AuthService, error) {
	return &AuthService{}, nil
}
