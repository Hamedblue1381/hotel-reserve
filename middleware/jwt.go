package middleware

import (
	"fmt"
	"hotel-project/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("-- JWT auth")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		fmt.Println("token no present in header")
		return fmt.Errorf("Unauthorized")
	}
	claims, err := parseToken(token[0])
	if err != nil {
		return fmt.Errorf("Unauthorized")
	}
	expires, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}
	if expires.Before(time.Now()) {
		return fmt.Errorf("Expired")
	}
	fmt.Println("claims: ", claims)
	return c.Next()
}

func parseToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signature algorithm", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token: ", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token: ", err)
		return nil, fmt.Errorf("unauthorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}

func CreateTokenFromUser(user *models.User) string {
	expires := time.Now().Add(time.Hour * 24)
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   jwt.NewNumericDate(expires),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
