package jwtservice

import (
	"authService/internal/cache"
	"authService/internal/config"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	TokenIssuer        = "wishlist_AuthService"
	RefreshTokenIssuer = "wishlist_AuthService_Refresh"
	CacheSetKeyValue   = "true"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.StandardClaims
}

type IJWTService interface {
	GenerateToken(userID, email, name string) (string, time.Time, error)
	GenerateRefreshToken(userID string) (string, time.Time, error)
	ValidateToken(token string) (*Claims, error)
	BlacklistToken(ctx context.Context, token string, expiration time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type JWTService struct {
	log   *slog.Logger
	cfg   *config.JWTConfig
	cache cache.ICache
}

func New(
	log *slog.Logger,
	cfg *config.JWTConfig,
	cache cache.ICache,

) *JWTService {
	return &JWTService{
		log:   log,
		cfg:   cfg,
		cache: cache,
	}
}

func (s *JWTService) GenerateToken(userID, email, name string) (string, time.Time, error) {
	expiresAt := time.Now().Add(s.cfg.TTL)

	claims := &Claims{
		UserID: userID,
		Email:  email,
		Name:   name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
			Issuer:    TokenIssuer,
			Subject:   userID,
			Id:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.Secret))

	return tokenString, expiresAt, err
}

func (s *JWTService) GenerateRefreshToken(userID string) (string, time.Time, error) {

	expiresAt := time.Now().Add(s.cfg.RefreshTTL)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    RefreshTokenIssuer,
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.Secret))

	return tokenString, expiresAt, err
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
func (s *JWTService) BlacklistToken(ctx context.Context, token string, expiration time.Duration) error {
	// Сохраняем токен в Redis до его истечения
	return s.cache.Set(ctx, "blacklist:"+token, CacheSetKeyValue, expiration)
}
func (s *JWTService) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := s.cache.Get(ctx, "blacklist:"+token)
	if err == cache.ErrorTokenNotSet {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return val == CacheSetKeyValue, nil
}
