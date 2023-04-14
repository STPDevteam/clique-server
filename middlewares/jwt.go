package middlewares

import (
	"fmt"
	"github.com/spf13/viper"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"stp_dao_v2/errs"

	"strings"
	"time"

	oo "github.com/Anna2024/liboo"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JWT struct {
	SigningKey []byte
}

type CustomClaims struct {
	Account string `json:"account"`
	UserId  int64  `json:"user_id"`
	jwt.StandardClaims
}

func NewJWT() *JWT {
	return &JWT{
		[]byte(viper.GetString("app.sign_message_pri_key")),
	}
}

// JWTAuth JWT is Gin framework middleware
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(consts.KEY_LOGIN, false)
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			return
		}
		authorizations := strings.Split(authorization, " ")
		if len(authorizations) != 2 || authorizations[0] != "Bearer" {
			return
		}
		token := authorizations[1]

		j := NewJWT()
		// parseToken parse token info
		claims, cErr := j.ParseToken(token)
		if cErr != nil {
			return
		}
		if claims.Subject != claims.Account {
			return
		}

		var err error
		var user db.TbAccountModel
		user, err = db.GetTbAccountModel(o.W("account", claims.Account))
		if err != nil {
			return
		}

		c.Set(consts.KEY_LOGIN, true)
		c.Set(consts.KEY_CURRENT_USER, user)
		c.Next()
	}
}

// JWTAuthForce JWT is Gin framework middleware, on this function, the jwt token is required.
func JWTAuthForce() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.Abort()
			c.Error(errs.ErrUnAuthorized)
			return
		}
		authorizations := strings.Split(authorization, " ")
		if len(authorizations) != 2 || authorizations[0] != "Bearer" {
			c.Abort()
			c.Error(errs.ErrUnAuthorized)
			return
		}
		token := authorizations[1]

		j := NewJWT()
		// parseToken parse token info
		claims, cErr := j.ParseToken(token)
		if cErr != nil {
			oo.LogW("%s: ParseToken: token:%s, error :%v", c.FullPath(), token, cErr)
			c.Abort()
			c.Error(cErr)
			return
		}
		if claims.Subject != claims.Account {
			oo.LogW("%s: ParseToken: token:%s, error : not consistent", c.FullPath(), token)
			c.Abort()
			c.Error(errs.ErrUnAuthorized)
			return
		}

		var err error
		var user db.TbAccountModel
		user, err = db.GetTbAccountModel(o.W("account", claims.Account))
		if err != nil {
			oo.LogW("%s: GetUserByAccount err, msg: %v", c.FullPath(), err)
			if err == oo.ErrNoRows {
				c.Abort()
				c.Error(errs.NewError(401, "User not exists."))
			} else {
				c.Abort()
				c.Error(errs.ErrServer)
			}
			return
		}

		c.Set(consts.KEY_LOGIN, true)
		c.Set(consts.KEY_CURRENT_USER, user)
		c.Next()
	}
}

// GenerateToken generates tokens used for auth
func (j *JWT) generateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateToken create the token
func CreateToken(account string, userId int64) (string, error) {
	j := NewJWT()
	claims := CustomClaims{
		account,
		userId,
		jwt.StandardClaims{
			Audience:  "clique",
			NotBefore: time.Now().Unix(),                                          // Signature effective time
			ExpiresAt: time.Now().Unix() + viper.GetInt64("app.jwt_expired_time"), // signature expiration time
			Issuer:    "clique",                                                   // signed issuer
			IssuedAt:  time.Now().Unix(),
			Subject:   account,
		},
	}
	token, err := j.generateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

// RefreshToken refresh the token
func (j *JWT) RefreshToken(tokenStr string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(time.Duration(viper.GetInt64("app.jwt_expired_time")) * time.Second).Unix()
		return j.generateToken(*claims)
	}
	return "", fmt.Errorf("couldn't handle this token")
}

// ParseToken parses token
func (j *JWT) ParseToken(tokenStr string) (*CustomClaims, *errs.CustomError) {
	tokenClaims, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errs.NewError(401, "That's not even a token.")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errs.NewError(401, "Login has expired.")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errs.NewError(401, "Token not active yet.")
			} else {
				return nil, errs.NewError(401, "Couldn't handle this token.")
			}
		}
	}

	if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, errs.NewError(401, "Couldn't handle this token.")
}
