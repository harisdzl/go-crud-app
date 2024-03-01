package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/implementations/auth"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)


func AuthHandler(p *base.Persistence) gin.HandlerFunc {

	return func(c *gin.Context) {
	   token := c.Request.Header.Get("Authorization")
	   
	   b := "Bearer "
	   if !strings.Contains(token, b) {
		  c.JSON(http.StatusForbidden, gin.H{"message": "Your request is not authorized", "status": entity.StatusError, "data": nil})
		  c.Abort()
		  return
	   }
	   t := strings.Split(token, b)
	   if len(t) < 2 {
		  c.JSON(http.StatusForbidden, gin.H{"message": "An authorization token was not supplied", "status": entity.StatusError, "data": nil})
		  c.Abort()
		  return
	   }
 
	   // Validate token	
	   ctx := c.Request.Context()

	   keyJWT := os.Getenv("JWT_SECRET")
	   v2, err2 := auth.NewAuthRepository(p, &ctx).ValidateToken(t[1], keyJWT)
 
	   if (err2 != nil) || (v2 != nil && !v2.Valid) {
		  c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization token", "status": entity.StatusError, "data": nil})
		  c.Abort()
		  return
	   }
 
	   //catch token
	   var tokenCatches jwt.Token
	   if v2.Valid {
		  tokenCatches = *v2
	   }
	   userIDInterface := tokenCatches.Claims.(jwt.MapClaims)["user_id"]

	   userID, ok := userIDInterface.(string)
	   if !ok {
		   fmt.Println("Invalid user ID format")
		   return
	   }

	   // Check redis blacklist tokens
	   isBlacklisted := auth.NewAuthRepository(p, &ctx).CheckBlacklistToken(t[1])
	   if userID == "" || isBlacklisted == nil {
		  c.JSON(http.StatusForbidden, gin.H{"message": "Invalid authorization token", "status": entity.StatusError, "data": nil})
		  c.Abort()
		  return
	   }
	   c.Set("userID", userID)
 
	   c.Next()
 
	}
}

