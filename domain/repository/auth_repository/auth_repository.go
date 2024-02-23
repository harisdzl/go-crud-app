package auth_repository

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
)
 
type AuthHandlerRepository interface {
	AuthenticateUser(username string, password string) (string, *customer_entity.Customer, error)
	BlacklistToken(token string) (error)
}

type AuthRepository interface {
	GenerateToken(key []byte, userId int64, credential string) (string, error)
	ValidateToken(tokenString string, key string) (*jwt.Token, error)
	GetCustomerWithUsername(username string) (*customer_entity.Customer, error)
	BlacklistToken(token string) (error)
}