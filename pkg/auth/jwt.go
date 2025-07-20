package auth

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token") // 标准错误定义
	secretKey       []byte                        // 包内私有的密钥存储
)

// 密钥通过 Init() 注入（推荐从环境变量读取）
// 存储为 []byte 类型（符合 JWT 库要求）
func Init(key string) {
	secretKey = []byte(key) // 初始化密钥（需在应用启动时调用）
}

// 令牌生成 (GenerateToken)
// JWT 组成：
// Header：自动生成（指定 HS256 算法）
// Payload：user_id：业务相关用户标识\exp：过期时间（RFC 7519 标准声明）\iat：签发时间（可选但推荐）
// 签名：使用 HMAC-SHA256 算法 + 密钥生成
// 返回值："头部.载荷.签名" 格式的字符串
func GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// 令牌解析 (ParseToken)
// 关键步骤：
// 签名验证：确保令牌未被篡改
// 算法检查：防止算法替换攻击
// 声明提取：从 payload 获取 user_id
// 类型转换：处理 JSON 数字到 Go 类型的映射

// 错误处理：
// 区分令牌无效和解析失败
// 始终返回标准化错误
func ParseToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secretKey, nil
	})

	if err != nil {
		return 0, err
	}
	// 类型断言提取声明
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, ErrInvalidToken
		}
		return uint(userID), nil
	}

	return 0, ErrInvalidToken
}

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//logger.Info("JWTMiddleware START", zap.String("token", c.GetHeader("Authorization")))
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header required"})
			return
		}

		userID, err := ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}
