package userservice

import (
	"authService/internal/domain/models"
	jwtservice "authService/internal/services/jwtService"
	"authService/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/asgwg01/wishlists/pkg/types/trace"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorInvalidLoginPwd  = errors.New("invalid login or password")
	ErrorUserNotExist     = errors.New("user not found")
	ErrorEmailExist       = errors.New("email already exists")
	ErrorTokenBlacklisted = errors.New("token is blacklisted")
)

type UserAuthService struct {
	log        *slog.Logger
	storage    storage.IStorage
	jwtService jwtservice.IJWTService
}

type UserAuth struct {
	User         models.User
	Token        string
	RefreshToken string
	ExpiresAt    time.Time
}

type IUserAuthService interface {
	CreateUser(ctx context.Context, email, pwd, name string) (UserAuth, error)
	Login(ctx context.Context, email, password string) (UserAuth, error)
	ValidateToken(ctx context.Context, token string) (models.User, error)
	RefreshToken(ctx context.Context, refreshToken string) (UserAuth, error)
	Logout(ctx context.Context, userID, token string) error
	UserInfo(ctx context.Context, userID string) (models.User, error)
}

func New(
	log *slog.Logger,
	storage storage.IStorage,
	jwt jwtservice.IJWTService,
) *UserAuthService {
	return &UserAuthService{
		log:        log,
		storage:    storage,
		jwtService: jwt,
	}
}

func (s *UserAuthService) CreateUser(ctx context.Context, email, pwd, name string) (UserAuth, error) {
	const logPrefix = "userService.CreateUser"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("email", email),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Create new user")

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error create pwd hash", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error create pwd hash %w", err)
	}

	newUser := models.User{
		Email:        email,
		PasswordHash: hash,
		Name:         name,
		CreateAt:     time.Now(),
		UpdateAt:     time.Now(),
	}
	newUser, err = s.storage.CreateUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, storage.ErrorEmailExist) {
			log.Error("Error create new user, email is exists", slog.String("err", err.Error()))
			return UserAuth{}, ErrorEmailExist

		}
		log.Error("Error create new user", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error create new user %w", err)
	}

	auth := UserAuth{
		User: newUser,
	}
	auth.Token, auth.ExpiresAt, err = s.jwtService.GenerateToken(newUser.ID.String(), newUser.Email, newUser.Name)
	if err != nil {
		log.Error("Error generate token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error generate token %w", err)
	}
	auth.RefreshToken, _, err = s.jwtService.GenerateRefreshToken(newUser.ID.String())
	if err != nil {
		log.Error("Error generate refresh token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error generate refresh token %w", err)
	}

	return auth, nil
}

func (s *UserAuthService) Login(ctx context.Context, email, pwd string) (UserAuth, error) {
	const logPrefix = "userService.Login"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("email", email),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Login")

	user, err := s.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotExist) {
			log.Error("Error get user, is not exist", slog.String("err", err.Error()))
			return UserAuth{}, ErrorInvalidLoginPwd

		}
		log.Error("Error get user", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error get user %w", err)
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pwd))
	if err != nil {
		log.Error("Error pwd comparing", slog.String("err", err.Error()))
		return UserAuth{}, ErrorInvalidLoginPwd
	}

	auth := UserAuth{
		User: user,
	}
	auth.Token, auth.ExpiresAt, err = s.jwtService.GenerateToken(user.ID.String(), user.Email, user.Name)
	if err != nil {
		log.Error("Error generate token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error generate token %w", err)
	}
	auth.RefreshToken, _, err = s.jwtService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		log.Error("Error generate refresh token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error generate refresh token %w", err)
	}

	return auth, nil
}

func (s *UserAuthService) ValidateToken(ctx context.Context, token string) (models.User, error) {
	const logPrefix = "userService.ValidateToken"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("ValidateToken")

	// Проверяем, не в черном ли списке токен
	isBlacklisted, err := s.jwtService.IsBlacklisted(ctx, token)
	if err != nil {
		log.Error("Error check black list", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("Error check black list %w", err)
	}
	if isBlacklisted {
		log.Error("Error user token in black list")
		return models.User{}, ErrorTokenBlacklisted
	}

	// Валидируем токен
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		log.Error("Error validate token", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("Error validate token %w", err)
	}

	uuid, err := uuid.Parse(claims.UserID)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("Error parse uuid %w", err)
	}

	user, err := s.storage.GetUserByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotExist) {
			log.Error("Error get user, is not exist", slog.String("err", err.Error()))
			return models.User{}, ErrorUserNotExist

		}
		log.Error("Error get user", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("Error get user %w", err)
	}

	return user, nil
}

func (s *UserAuthService) RefreshToken(ctx context.Context, refreshToken string) (UserAuth, error) {
	const logPrefix = "userService.RefreshToken"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("RefreshToken")

	// Проверяем, не в черном ли списке токен
	isBlacklisted, err := s.jwtService.IsBlacklisted(ctx, refreshToken)
	if err != nil {
		log.Error("Error check black list", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error check black list %w", err)
	}
	if isBlacklisted {
		log.Error("Error user token in black list")
		return UserAuth{}, ErrorTokenBlacklisted
	}

	// Валидируем токен
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		log.Error("Error validate token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("error validate token %w", err)
	}

	// Проверяем, что это действительно refresh token (по issuer или типу)
	if claims.Issuer != jwtservice.RefreshTokenIssuer {
		log.Error("Token is not a refresh token")
		return UserAuth{}, fmt.Errorf("token is not a refresh token %w", err)
	}

	uuid, err := uuid.Parse(claims.UserID)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("error parse uuid %w", err)
	}

	user, err := s.storage.GetUserByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotExist) {
			log.Error("Error get user, is not exist", slog.String("err", err.Error()))
			return UserAuth{}, ErrorUserNotExist

		}
		log.Error("Error get user", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error get user %w", err)
	}

	// Добавляем старый refresh token в черный список (security best practice)
	// Это предотвращает повторное использование того же refresh token
	expiresIn := time.Until(time.Unix(0, claims.ExpiresAt))
	if err := s.jwtService.BlacklistToken(ctx, refreshToken, expiresIn); err != nil {
		log.Error("Error move to black list token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error move to black list token %w", err)
	}

	// Генерируем новую пару токенов
	auth := UserAuth{
		User: user,
	}
	auth.Token, auth.ExpiresAt, err = s.jwtService.GenerateToken(user.ID.String(), user.Email, user.Name)
	if err != nil {
		log.Error("Error generate token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error generate token %w", err)
	}
	auth.RefreshToken, _, err = s.jwtService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		log.Error("Error generate refresh token", slog.String("err", err.Error()))
		return UserAuth{}, fmt.Errorf("Error generate refresh token %w", err)
	}

	return auth, nil
}

func (s *UserAuthService) Logout(ctx context.Context, userID, token string) error {
	const logPrefix = "userService.Logout"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("email", userID),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("Logout")

	// Получаем claims токена чтобы узнать время истечения
	claims, err := s.jwtService.ValidateToken(token)
	if err != nil {
		log.Error("Error validate token", slog.String("err", err.Error()))
		return fmt.Errorf("error validate token %w", err)
	}

	// Добавляем токен в черный список до его истечения
	expiration := time.Until(time.Unix(0, claims.ExpiresAt))
	err = s.jwtService.BlacklistToken(ctx, token, expiration)
	if err != nil {
		log.Error("Error move to black list token", slog.String("err", err.Error()))
		return fmt.Errorf("Error move to black list token %w", err)
	}

	return nil
}
func (s *UserAuthService) UserInfo(ctx context.Context, userID string) (models.User, error) {
	const logPrefix = "userService.UserInfo"
	log := s.log.With(
		slog.String("where", logPrefix),
		slog.String("userID", userID),
		slog.String("trace_id", trace.GetTraceID(ctx)),
	)

	log.Info("UserInfo")

	uuid, err := uuid.Parse(userID)
	if err != nil {
		log.Error("Error parse uuid", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("error parse uuid %w", err)
	}

	user, err := s.storage.GetUserByID(ctx, uuid)
	if err != nil {
		if errors.Is(err, storage.ErrorUserNotExist) {
			log.Error("Error get user, is not exist", slog.String("err", err.Error()))
			return models.User{}, ErrorUserNotExist

		}
		log.Error("Error get user", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("Error get user %w", err)
	}

	return user, nil
}
