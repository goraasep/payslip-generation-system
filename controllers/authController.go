package controllers

import (
	"os"

	"github.com/goraasep/payslip-generation-system/config"
	"github.com/goraasep/payslip-generation-system/dto"
	"github.com/goraasep/payslip-generation-system/models"
	"github.com/goraasep/payslip-generation-system/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// POST /register
func Register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		response.InternalError(c, "Failed to hash password")
		return
	}

	var role models.Role
	if err := config.DB.Where("name = ?", "USER").First(&role).Error; err != nil {
		response.InternalError(c, "Role USER not found")
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
		Roles:    []*models.Role{&role},
	}

	if err := config.DB.Create(&user).Error; err != nil {
		response.InternalError(c, "Email may already exist")
		return
	}
	response.Success(c, "User registered successfully", user)
}

// POST /login
func Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID)
	if err != nil {
		response.InternalError(c, "Failed to generate token")
		return
	}

	tokenResponse := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.Success(c, "Login success", tokenResponse)
}

// POST /refresh
func Refresh(c *gin.Context) {
	var input dto.RefrehRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, "Missing refresh token")
		return
	}

	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		response.Unauthorized(c, "Invalid refresh token")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		response.InternalError(c, "Invalid token claims")
		return
	}

	userID := uint(claims["user_id"].(float64))
	accessToken, refreshToken, err := utils.GenerateTokens(userID)
	if err != nil {
		response.InternalError(c, "Failed to generate token")
		return
	}

	tokenResponse := dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.Success(c, "Token refreshed", tokenResponse)
}

// GET /api/me (protected)
func Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		response.NotFound(c, "User not found")
		return
	}
	response.Success(c, "User found", user)
}
