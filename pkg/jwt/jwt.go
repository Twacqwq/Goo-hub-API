// Package jwt 处理 JWT 认证
package jwt

import (
	"errors"
	"strings"
	"thub/pkg/config"
	"thub/pkg/logger"
	"time"

	"thub/pkg/app"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌格式有误")
	ErrTokenInvalid           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("需要认证才能访问")
	ErrHeaderMalformed        error = errors.New("请求头中 Authorization 格式有误")
)

// 定义一个JWT对象
type JWT struct {
	// 密钥, 加密JWT, 读取配置 app.key
	SignKey []byte

	// 刷新Token的最大过期时间
	MaxRefresh time.Duration
}

// 自定义载荷Claims
type JWTCustomClaims struct {
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	ExpireAtTime int64  `json:"expire_time"`

	// StandardClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwtpkg.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		SignKey:    []byte(config.GetString("app.key")),
		MaxRefresh: time.Duration(config.GetInt64("jwt.max_refresh_time")) * time.Minute,
	}
}

// 解析 Token, 中间件调用
func (jwt *JWT) ParseToken(c *gin.Context) (*JWTCustomClaims, error) {
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}
	// 调用jwt库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)
	// 解析出错
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	// 将token中的claims信息解析出来和 JWTCustomClaims结构体校验
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// 更新Token 用以 refresh token 接口
func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	token, err := jwt.parseTokenString(tokenString)
	if err != nil {
		validation, ok := err.(*jwtpkg.ValidationError)
		// 符合 Refresh条件 但 单一的报错
		if !ok || validation.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}

	// 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*JWTCustomClaims)

	// 检查是否过了 最大允许刷新时间
	x := app.TimenowInTimezone().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt > x {
		claims.StandardClaims.ExpiresAt = jwt.expireAtTime()
		return jwt.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

// 生成Token
func (jwt *JWT) IssueToken(userID, userName string) string {
	// 构造用户 claims 信息(载荷)
	expireAtTime := jwt.expireAtTime()
	claims := JWTCustomClaims{
		userID,
		userName,
		expireAtTime,
		jwtpkg.StandardClaims{
			NotBefore: app.TimenowInTimezone().Unix(), // 签名生效时间
			IssuedAt:  app.TimenowInTimezone().Unix(), // 首次签名时间(后续刷新token不会更新)
			ExpiresAt: expireAtTime,                   // 签名过期时间
			Issuer:    config.GetString("app.name"),   // 颁发者
		},
	}
	// 根据claims生成token
	token, err := jwt.createToken(claims)
	if err != nil {
		logger.LogIf(err)
		return ""
	}

	return token
}

// 内部创建token
func (jwt *JWT) createToken(claims JWTCustomClaims) (string, error) {
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.SignKey)
}

// 过期时间
func (jwt *JWT) expireAtTime() int64 {
	timenow := app.TimenowInTimezone()

	var expireTime int64
	if config.GetBool("app.debug") {
		expireTime = config.GetInt64("jwt.debug_expire_time")
	} else {
		expireTime = config.GetInt64("jwt.expire_time")
	}

	expire := time.Duration(expireTime) * time.Minute
	return timenow.Add(expire).Unix()
}

// 解析token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.SignKey, nil
	})
}

// 从Authorization中解析token
// Authorization:Bearer xxxxxx
func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}

	return parts[1], nil
}
