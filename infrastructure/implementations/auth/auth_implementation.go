package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/customer_entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/cache"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)



type AuthRepo struct {
	p *base.Persistence
	c *context.Context
}


func NewAuthRepository(p *base.Persistence, c *context.Context) *AuthRepo {
	return &AuthRepo{p, c}
}



func (a AuthRepo) GenerateToken(key []byte, userId int64, credential string) (string, error) {

	//new token
	token := jwt.New(jwt.SigningMethodHS256)
 
	// Claims
	claims := make(jwt.MapClaims)
	userIdString := fmt.Sprint(userId)
	claims["user_id"] = userIdString
	claims["credential"] = credential
	claims["exp"] = time.Now().Add(time.Hour*720).UnixNano() / int64(time.Millisecond)
 
	//Set user roles
	//claims["roles"] = roles
 
	token.Claims = claims
 
	// Sign and get as a string
	tokenString, err := token.SignedString(key)
	return tokenString, err
}

func (a AuthRepo) ValidateToken(tokenString string, key string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	   return []byte(key), nil
	})
 
	return token, err
}

func (a AuthRepo) GetCustomerWithUsername(username string) (*customer_entity.Customer, error) {
	var customer *customer_entity.Customer
	err := a.p.DB.Debug().Where("username = ?", username).Take(&customer).Error

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (a AuthRepo) BlacklistToken(token string) error {
	cacheRepo := cache.NewCacheRepository("Redis", a.p)

	err := cacheRepo.SetKey(fmt.Sprintf("%v_JWTTOKEN", token), token, time.Hour * 720)
	if err != nil {
		return err
	}
	return nil
}

func (a AuthRepo) CheckBlacklistToken(token string) error {
	cacheRepo := cache.NewCacheRepository("Redis", a.p)

	err := cacheRepo.GetKey(fmt.Sprintf("%v_JWTTOKEN", token), token)

	if err != nil {
		return err
	}


	return nil
}