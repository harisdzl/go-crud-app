package application

import (
	"context"
	"errors"
	"os"

	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/auth_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/auth"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/utils/security"
)

type AuthApp struct {
	p *base.Persistence
	c *context.Context
}

func NewAuthApplication(p *base.Persistence, c *context.Context) auth_repository.AuthHandlerRepository {
	return &AuthApp{p, c}
}

func (a *AuthApp) AuthenticateUser(username string, password string) (string, *customer_entity.Customer, error) {
	authRepo := auth.NewAuthRepository(a.p, a.c)
	customer, err := authRepo.GetCustomerWithUsername(username)

	if err != nil {
		return "", nil, err
	}

	verifiedPassword := security.VerifyPassword(customer.Password, password)

	if verifiedPassword != nil {
		return "", nil, verifiedPassword
	}

	verifiedUsername := customer.Username == username

	if !verifiedUsername {
		return "", nil, errors.New("Invalid username")
	}

	token, err := authRepo.GenerateToken([]byte(os.Getenv("JWT_SECRET")), int64(customer.ID), customer.Username)

	if err != nil {
		return "", nil, err
	}

	return token, customer, nil
}

func (a *AuthApp) BlacklistToken(token string) error {
	authRepo := auth.NewAuthRepository(a.p, a.c)
	return authRepo.BlacklistToken(token)
}
