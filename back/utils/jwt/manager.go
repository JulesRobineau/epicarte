package jwt

import (
	"fmt"
	"gin-template/config"
	"gin-template/pkg/model/enum"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtManager struct {
	Secret    string `json:"secret"`
	ExpiresIn int    `json:"expires_in"`
}

type JwtToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    time.Time `json:"expires_in" format:"unix"`
	TokenID      string    `json:"-"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserId uint64    `json:"user_id"`
	Role   enum.Role `json:"role"`
}

func NewJwtManager(secret string, expiresIn int) JwtManager {
	return JwtManager{
		Secret:    secret,
		ExpiresIn: expiresIn,
	}
}

func GenerateTokens(userId uint64, role enum.Role, conf config.JwtConfig) (JwtToken, error) {
	tokenId := fmt.Sprintf("%d:%d", userId, time.Now().Unix())
	now := time.Now().UTC().Add(time.Duration(conf.Expiration) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now),
			ID:        tokenId,
		},
		UserId: userId,
		Role:   role,
	})

	refresh, err := token.SignedString([]byte(conf.Secret))
	if err != nil {
		return JwtToken{}, err
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID: tokenId,
		},
		UserId: userId,
		Role:   role,
	})
	access, err := token.SignedString([]byte(conf.Secret))
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    now,
		TokenID:      tokenId,
	}, nil
}

func ParseToken(token, secret string) (*Claims, error) {
	cl := &Claims{}
	_, err := jwt.ParseWithClaims(token, cl, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	return cl, nil
}

func (j *JwtToken) RefreshTokenWithToken(refreshToken string, conf config.JwtConfig) error {
	var cl *Claims
	var err error
	if refreshToken == "" {
		cl = &Claims{
			Role: enum.SUPERADMIN,
		}
	} else {
		cl, err = ParseToken(refreshToken, conf.Secret)
	}
	if err != nil {
		return err
	}

	now := time.Now().UTC().Add(time.Duration(conf.Expiration) * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now),
		},
		UserId: cl.UserId,
		Role:   cl.Role,
	})

	access, err := token.SignedString([]byte(conf.Secret))
	if err != nil {
		return err
	}

	j.AccessToken = access
	j.ExpiresIn = now
	return nil
}
