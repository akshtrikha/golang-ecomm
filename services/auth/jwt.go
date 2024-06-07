package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/akshtrikha/golang-ecomm/config"
	"github.com/golang-jwt/jwt"
)

// GenerateJWT function generates and returns the jwt token for the key provided
func GenerateJWT(secret string, userID int) (string, error) {
    expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userID": strconv.Itoa(userID),
        "expiredAt": time.Now().Add(expiration).Unix(),
    })

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken Function
func VerifyToken(tokenString string, secretKey string) (*jwt.Token, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Check the signing method to make sure it's what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

        // Return your secret key:
        return []byte(secretKey), nil
    })
    return token, err
}

// VerifyTokenAndClaims function
func VerifyTokenAndClaims(tokenString string, secretKey string) (jwt.MapClaims, error) {
    token, err := VerifyToken(tokenString, secretKey)
    if err != nil {
        return nil, err
    }

    // Check if token is valid
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil, fmt.Errorf("invalid claims format")
    }

    // Validate expiration time
    exp, ok := claims["expiredAt"].(float64)
    if !ok {
        return nil, fmt.Errorf("missing expiration time claim")
    }
    if int64(exp) < time.Now().Unix() {
        return nil, fmt.Errorf("token expired")
    }

    // ... other custom claim validations ...

    return claims, nil
}
