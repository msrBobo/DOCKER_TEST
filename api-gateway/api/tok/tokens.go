package tokens

// import (
// 	"DOCKER_TEST/api-gateway/api/models"
// 	"DOCKER_TEST/api-gateway/config"
// 	"fmt"
// 	"time"

// 	"github.com/golang-jwt/jwt"
// )

// func GenerateTokens(user models.User) (string, string, error) {
// 	c := config.Config{}

// 	// Generate access token
// 	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userID":   user.Id,
// 		"userRole": user.Role,
// 		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Short expiration time
// 	})

// 	accessTokenString, err := accessToken.SignedString([]byte(c.JWT_SECRET_KEY))
// 	if err != nil {
// 		return "", "", err
// 	}

// 	// Generate refresh token
// 	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"userID":   user.Id,
// 		"userRole": user.Role,
// 		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Longer expiration time
// 	})
// 	refreshTokenString, err := refreshToken.SignedString([]byte(c.JWT_SECRET_KEY))
// 	if err != nil {
// 		return "", "", err
// 	}

// 	return accessTokenString, refreshTokenString, nil
// }

// func ParseRefreshToken(tokenString string) (string, error) {
// 	c := config.Config{}
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(c.JWT_SECRET_KEY), nil
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		userID := claims["userID"].(string)
// 		return userID, nil
// 	}
// 	return "", fmt.Errorf("invalid token")
// }
