package pkg

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"uid"`
	jwt.RegisteredClaims
}

// GenerateJWTToken 生成 JWT 令牌
func GenerateJWTToken(userID uint, jwtSecret []byte, expiresIn time.Duration) (string, error) {
	now := time.Now()

	if expiresIn <= 0 {
		expiresIn = time.Hour * 24
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		// "sub": userID,                    // 用户ID
		// "exp": now.Add(expiresIn).Unix(), // 过期时间:1天
		// "iat": now.Unix(),                // 签发时间
		// "nbf": now.Unix(),                // 生效时间
		// "jti": uuid.New().String(),       // 唯一标识
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),                // 签发时间
			NotBefore: jwt.NewNumericDate(now),                // 生效时间
			ID:        uuid.New().String(),                    // 唯一标识
		},
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}
	return tokenString, nil

}

// ParseJWTToken 解析 JWT 令牌
// 参数: tokenString 令牌字符串, jwtSecret 密钥
//
// 返回: 用户ID, 错误
//
// 错误: 过期, 签名错误, 格式错误, token 无效
func ParseJWTToken(tokenString string, jwtSecret []byte) (uint, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		// 验证签名方法
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	// 过期，签名错误，格式错误，token 无效
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}
