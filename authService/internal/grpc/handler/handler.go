package handler

import (
	userservice "authService/internal/services/userService"
	"context"
	"errors"

	authv1 "github.com/asgwg01/wishlists/pkg/proto/auth/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrorEmailExist = errors.New("email already exist")
)

// GrpcHandler реализует сгенерированный интерфейс
type GrpcHandler struct {
	authv1.UnimplementedAuthServiceServer
	userService userservice.UserAuthService
}

func NewGrpcHandler(userService userservice.UserAuthService) *GrpcHandler {
	return &GrpcHandler{
		userService: userService,
	}
}

// Регистрация нового пользователя
func (h *GrpcHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.AuthResponse, error) {
	err := validateRegisterRequest(req)
	if err != nil {
		return nil, err
	}

	userAuth, err := h.userService.CreateUser(ctx, req.GetEmail(), req.GetPassword(), req.GetName())
	if err != nil {
		switch err {
		case userservice.ErrorEmailExist:
			return nil, status.Error(codes.AlreadyExists, "email already registered")
		default:
			return nil, status.Error(codes.Internal, "failed to register user")
		}
	}

	return &authv1.AuthResponse{
		Token:        userAuth.Token,
		RefreshToken: userAuth.RefreshToken,
		UserId:       userAuth.User.ID.String(),
		ExpiresAt:    userAuth.ExpiresAt.Unix(),
	}, nil
}

// Вход
func (h *GrpcHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.AuthResponse, error) {
	err := validateLoginRequest(req)
	if err != nil {
		return nil, err
	}

	userAuth, err := h.userService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		switch err {
		case userservice.ErrorInvalidLoginPwd:
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		default:
			return nil, status.Error(codes.Internal, "failed to login user")
		}
	}

	return &authv1.AuthResponse{
		Token:        userAuth.Token,
		RefreshToken: userAuth.RefreshToken,
		UserId:       userAuth.User.ID.String(),
		ExpiresAt:    userAuth.ExpiresAt.Unix(),
	}, nil
}

// Выход
func (h *GrpcHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponce, error) {
	err := validateLogoutRequest(req)
	if err != nil {
		return nil, err
	}

	err = h.userService.Logout(ctx, req.GetUserId(), req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to logout user")
	}

	return &authv1.LogoutResponce{
		Sucess: true,
	}, nil
}

// Валидация / чек корректности токена
func (h *GrpcHandler) Validate(ctx context.Context, req *authv1.ValidateRequest) (*authv1.UserInfoResponce, error) {
	err := validateValidateRequest(req)
	if err != nil {
		return nil, err
	}

	userInfo, err := h.userService.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to validate token")
	}

	return &authv1.UserInfoResponce{
		UserId: userInfo.ID.String(),
		Email:  userInfo.Email,
		Name:   userInfo.Name,
	}, nil
}

// обновление токена
func (h *GrpcHandler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.AuthResponse, error) {
	err := validateRefreshRequest(req)
	if err != nil {
		return nil, err
	}

	userAuth, err := h.userService.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to refresh token")
	}

	return &authv1.AuthResponse{
		Token:        userAuth.Token,
		RefreshToken: userAuth.RefreshToken,
		UserId:       userAuth.User.ID.String(),
		ExpiresAt:    userAuth.ExpiresAt.Unix(),
	}, nil
}

func (h *GrpcHandler) GetUserInfo(ctx context.Context, req *authv1.UserInfoRequest) (*authv1.UserInfoResponce, error) {
	err := validateUserInfoRequest(req)
	if err != nil {
		return nil, err
	}

	userInfo, err := h.userService.UserInfo(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, userservice.ErrorUserNotExist) {
			return nil, status.Error(codes.NotFound, "user is not exist")
		}
		return nil, status.Error(codes.Internal, "failed to get user info")
	}

	return &authv1.UserInfoResponce{
		UserId: userInfo.ID.String(),
		Email:  userInfo.Email,
		Name:   userInfo.Name,
	}, nil
}

func validateRegisterRequest(req *authv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}

	return nil
	
}

func validateLoginRequest(req *authv1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func validateLogoutRequest(req *authv1.LogoutRequest) error {
	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetToken() == "" {
		return status.Error(codes.InvalidArgument, "token is required")
	}

	return nil
}

func validateValidateRequest(req *authv1.ValidateRequest) error {
	if req.GetToken() == "" {
		return status.Error(codes.InvalidArgument, "token is required")
	}

	return nil
}

func validateRefreshRequest(req *authv1.RefreshRequest) error {
	if req.GetRefreshToken() == "" {
		return status.Error(codes.InvalidArgument, "refresh_token is required")
	}

	return nil
}
func validateUserInfoRequest(req *authv1.UserInfoRequest) error {
	if req.GetUserId() == "" {
		return status.Error(codes.InvalidArgument, "user_id is required")
	}

	return nil
}
