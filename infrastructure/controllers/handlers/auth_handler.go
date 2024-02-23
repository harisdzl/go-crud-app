package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/application"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/auth_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type Auth struct {
	AuthRepo    auth_repository.AuthHandlerRepository
	Persistence *base.Persistence
}

type LoginDetails struct {
    Username   string   `json:"username"`
    Password string `json:"password"`
}

// NewAuth initializes a new Auth instance.
func NewAuth(p *base.Persistence) *Auth {
	return &Auth{
		Persistence: p,
	}
}

//	@Summary		User login
//	@Description	Authenticates user credentials and returns an access token upon successful authentication.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Login Details	body		LoginDetails	true	"User credentials"
//	@Success		200			{object}	entity.ResponseContext		"Success"
//	@Failure		400			{object}	entity.ResponseContext		"Bad request"
//	@Failure		401			{object}	entity.ResponseContext		"Unauthorized"
//	@Router			/login [post]
func (au *Auth) Login(c *gin.Context) {
	loginDetails := LoginDetails{}
	responseContextData := entity.ResponseContext{Ctx: c}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest,
			responseContextData.ResponseData(entity.StatusFail, err.Error(), ""))
		return
	}

	au.AuthRepo = application.NewAuthApplication(au.Persistence)

	token, customer, tokenErr := au.AuthRepo.AuthenticateUser(loginDetails.Username, loginDetails.Password)

	if tokenErr != nil {
		c.JSON(http.StatusUnauthorized,
			responseContextData.ResponseData(entity.StatusFail, tokenErr.Error(), ""))
		return
	}

	c.Header("Authorization", "Bearer "+token) // Set Authorization header
	userData := make(map[string]interface{})
	userData["access_token"] = token
	userData["id"] = customer.ID
	userData["username"] = customer.Username
	userData["name"] = customer.Name

	c.JSON(http.StatusOK, responseContextData.ResponseData(entity.StatusSuccess, "Token sent.", userData))
}

//	@Summary		User logout
//	@Description	Logs out the user by blacklisting the provided access token on Redis
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Access Token"
//	@Success		200				{object}	entity.ResponseContext	"Success"
//	@Failure		400				{object}	entity.ResponseContext	"Bad request"
//	@Failure		500				{object}	entity.ResponseContext	"Internal server error"
//	@Router			/logout [post]
func (au *Auth) Logout(c *gin.Context) {
	responseContextData := entity.ResponseContext{Ctx: c}
	authHeader := c.Request.Header.Get("Authorization")

	token := strings.TrimPrefix(authHeader, "Bearer ")
	accessToken := map[string]string{
		"access_token": token,
	}

	if accessToken["access_token"] == "" {
		c.JSON(http.StatusBadRequest,
			responseContextData.ResponseData(entity.StatusFail, "No token found", ""))
		return
	}

	// Add the access token to the blacklist in Redis
	au.AuthRepo = application.NewAuthApplication(au.Persistence)
	if err := au.AuthRepo.BlacklistToken(accessToken["access_token"]); err != nil {
		c.JSON(http.StatusInternalServerError,
			responseContextData.ResponseData(entity.StatusFail, "Failed to logout", nil))
		return
	}

	c.JSON(http.StatusOK,
		responseContextData.ResponseData(entity.StatusSuccess, "Logged out successfully", nil))
}
