package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nafisalfiani/ketson-go-lib/codes"
	"github.com/nafisalfiani/ketson-go-lib/errors"
	"github.com/nafisalfiani/ketson-go-lib/header"
	"github.com/nafisalfiani/ketson-go-lib/log"
	"github.com/nafisalfiani/ketson-go-lib/parser"
)

type contextKey string

const (
	userAuthInfo contextKey = "UserAuthInfo"
)

type Interface interface {
	VerifyToken(ctx context.Context, tokenId string) (User, error)
	CreateToken(ctx context.Context, user User) (Token, error)

	SetUserAuthInfo(ctx context.Context, user User, token *Token) context.Context
	GetUserAuthInfo(ctx context.Context) (UserAuthInfo, error)
}

type auth struct {
	log        log.Interface
	json       parser.JSONInterface
	httpClient *http.Client
	conf       Config
}

type Config struct {
	AccessTokenKey      string        `env:"ACCESS_TOKEN_KEY"`
	AccessTokenDuration time.Duration `env:"ACCESS_TOKEN_DURATION"`
}

func Init(cfg Config, log log.Interface, json parser.JSONInterface, httpClient *http.Client) Interface {
	return &auth{
		log:        log,
		json:       json,
		httpClient: httpClient,
		conf:       cfg,
	}
}

func (a *auth) VerifyToken(ctx context.Context, accessToken string) (User, error) {
	var user User
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewWithCode(codes.CodeAuth, "Unexpected signing method")
		}
		return []byte(a.conf.AccessTokenKey), nil
	})
	if err != nil {
		return user, errors.NewWithCode(codes.CodeAuth, err.Error())
	}

	if !token.Valid {
		return user, errors.NewWithCode(codes.CodeAuthInvalidToken, "invalid token")
	}

	user, expiry, err := a.extractClaims(token.Claims.(jwt.MapClaims))
	if err != nil {
		return user, err
	}

	if time.Unix(expiry, 0).Before(time.Now()) {
		return user, errors.NewWithCode(codes.CodeAuthAccessTokenExpired, "access token expired")
	}

	return user, nil
}

func (a *auth) CreateToken(ctx context.Context, user User) (Token, error) {
	accessExpiryTime := time.Now().Add(a.conf.AccessTokenDuration).Unix()
	accessClaims := jwt.MapClaims{
		"user:id":       user.Id,
		"user:name":     user.Name,
		"user:email":    user.Email,
		"user:username": user.Username,
		"user:role":     user.Role,
		"user:scopes":   user.Scopes,
		"expiry":        accessExpiryTime,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(a.conf.AccessTokenKey))
	if err != nil {
		return Token{}, err
	}

	token := Token{
		TokenType:       header.AuthorizationBearer,
		AccessToken:     accessTokenString,
		AccessExpiresIn: accessExpiryTime,
	}

	return token, nil
}

func (a *auth) extractClaims(claims jwt.MapClaims) (User, int64, error) {
	var user User
	var expiredIn int64
	var ok bool

	user.Id, ok = claims["user:id"].(string)
	if !ok {
		return user, expiredIn, errors.NewWithCode(codes.CodeAuth, "invalid id format")
	}

	user.Name, ok = claims["user:name"].(string)
	if !ok {
		return user, expiredIn, errors.NewWithCode(codes.CodeAuth, "invalid name format")
	}

	user.Email, ok = claims["user:email"].(string)
	if !ok {
		return user, expiredIn, errors.NewWithCode(codes.CodeAuth, "invalid email format")
	}

	user.Username, ok = claims["user:username"].(string)
	if !ok {
		return user, expiredIn, errors.NewWithCode(codes.CodeAuth, "invalid username format")
	}

	user.Role, ok = claims["user:role"].(string)
	if !ok {
		return user, expiredIn, errors.NewWithCode(codes.CodeAuth, "invalid role format")
	}

	// scopes, ok := claims["user:scopes"].([]any)
	// if !ok {
	// 	return user, expiredIn, errors.NewWithCode(codes.CodeAuth, "Invalid scopes format")
	// }

	// for i := range scopes {
	// 	if s, ok := scopes[i].(string); ok {
	// 		user.Scopes = append(user.Scopes, s)
	// 	}
	// }

	expiry, ok := claims["expiry"].(float64)
	if !ok {
		return user, expiredIn, errors.NewWithCode(codes.CodeAuth, "invalid expiry format")
	}

	return user, int64(expiry), nil
}

func (a *auth) SetUserAuthInfo(ctx context.Context, user User, token *Token) context.Context {
	userauth := UserAuthInfo{
		User:  user,
		Token: token,
	}

	return context.WithValue(ctx, userAuthInfo, userauth)
}

func (a *auth) GetUserAuthInfo(ctx context.Context) (UserAuthInfo, error) {
	user, ok := ctx.Value(userAuthInfo).(UserAuthInfo)
	if !ok {
		return user, errors.NewWithCode(codes.CodeAuthFailure, "failed getting user auth info")
	}

	return user, nil
}
