package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/helpers"
	"github.com/goraasep/payslip-generation-system/models"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helpers.ResponseHelper.Unauthorized(c, "Missing token")
			c.Abort()
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})
		if err != nil || !token.Valid {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			helpers.ResponseHelper.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			helpers.ResponseHelper.Unauthorized(c, "Invalid token claims")
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			helpers.ResponseHelper.Unauthorized(c, "Invalid token payload")
			c.Abort()
			return
		}

		c.Set("user_id", uint(userIDFloat))
		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			helpers.ResponseHelper.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}

		var user models.User
		if err := config.DB.Preload("Roles").First(&user, userID).Error; err != nil {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			helpers.ResponseHelper.Unauthorized(c, "User not found")
			c.Abort()
			return
		}

		userRoles := map[string]bool{}
		for _, role := range user.Roles {
			userRoles[role.Name] = true
		}

		for _, required := range roles {
			if userRoles[required] {
				c.Next()
				return
			}
		}

		// c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden – insufficient role"})
		helpers.ResponseHelper.Unauthorized(c, "Forbidden – insufficient role")
		c.Abort()
	}
}
