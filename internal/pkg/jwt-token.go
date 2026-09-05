package pkg

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWTToken 生成 JWT 令牌
func GenerateJWTToken(userID uint, jwtSecret []byte, expiresIn time.Duration) (string, error) {
	now := time.Now()

	if expiresIn <= 0 {
		expiresIn = time.Hour * 24
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,                    // 用户ID
		"exp": now.Add(expiresIn).Unix(), // 过期时间:1天
		"iat": now.Unix(),                // 签发时间
		"nbf": now.Unix(),                // 生效时间
		"jti": uuid.New().String(),       // 唯一标识
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil

}
