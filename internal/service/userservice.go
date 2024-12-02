package service

import (
	"context"
	"errors"
	pb "github.com/TiktokCommence/userService/api/user/v1"
	"github.com/TiktokCommence/userService/internal/errcode"
	"github.com/TiktokCommence/userService/internal/tool"
)

type UserServiceService struct {
	pb.UnimplementedUserServiceServer
	userHandler UserHandler
}

func NewUserServiceService(userHandler UserHandler) *UserServiceService {
	return &UserServiceService{
		userHandler: userHandler,
	}
}

func (s *UserServiceService) Register(ctx context.Context, req *pb.RegisterReq) (*pb.RegisterResp, error) {
	if req.GetPassword() != req.GetConfirmPassword() {
		return &pb.RegisterResp{}, ErrPasswordsDoNotMatch
	}
	if !tool.CheckPassword(req.GetPassword()) {
		return &pb.RegisterResp{}, ErrPasswordNotValid
	}
	if !s.userHandler.VerifyCode(ctx, req.GetEmail(), req.GetVerifyCode()) {
		return &pb.RegisterResp{}, ErrEmailVerifyCode
	}
	userID, err := s.userHandler.CreateUser(ctx, req.GetEmail(), req.GetPassword())
	if errors.Is(err, errcode.UserAlreadyExists) {
		return &pb.RegisterResp{}, ErrUserAlreadyExists
	}
	if err != nil {
		return &pb.RegisterResp{}, ErrCreateUser
	}
	return &pb.RegisterResp{
		UserId: userID,
	}, nil
}
func (s *UserServiceService) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	user, err := s.userHandler.GetUserInfoByEmail(ctx, req.GetEmail())
	if errors.Is(err, errcode.UserNotFound) {
		return &pb.LoginResp{}, ErrUserNotFound
	}
	if err != nil {
		return &pb.LoginResp{}, ErrLogin
	}
	if user.Password != req.GetPassword() {
		return &pb.LoginResp{}, ErrPasswordIncorrect
	}
	return &pb.LoginResp{
		UserId: user.ID,
	}, nil
}
func (s *UserServiceService) Logout(ctx context.Context, req *pb.LogoutReq) (*pb.LogoutResp, error) {
	err := s.userHandler.Logout(ctx, req.GetUserId())
	if errors.Is(err, errcode.UserNotFound) {
		return &pb.LogoutResp{Success: false}, ErrUserNotFound
	}
	if err != nil {
		return &pb.LogoutResp{Success: false}, ErrLogout
	}
	return &pb.LogoutResp{Success: true}, nil
}
func (s *UserServiceService) DeleteUser(ctx context.Context, req *pb.DeleteReq) (*pb.DeleteResp, error) {
	err := s.userHandler.DeleteUser(ctx, req.GetUserId())
	if errors.Is(err, errcode.UserNotFound) {
		return &pb.DeleteResp{Success: false}, ErrUserNotFound
	}
	if err != nil {
		return &pb.DeleteResp{Success: false}, ErrDeleteUser
	}
	return &pb.DeleteResp{Success: true}, nil
}
func (s *UserServiceService) UpdateUser(ctx context.Context, req *pb.UpdateReq) (*pb.UpdateResp, error) {
	user, err := s.userHandler.GetUserInfoByID(ctx, req.GetId())
	if errors.Is(err, errcode.UserNotFound) {
		return &pb.UpdateResp{Success: false}, ErrUserNotFound
	}
	if err != nil {
		return &pb.UpdateResp{Success: false}, ErrUpdateUser
	}
	if req.Name != nil {
		user.Name = req.Name
	}
	if req.Age != nil {
		user.Age = req.Age
	}
	if req.Addr1 != nil {
		user.Addr1 = req.Addr1
	}
	if req.Addr2 != nil {
		user.Addr2 = req.Addr2
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	err = s.userHandler.UpdateUserInfo(ctx, user)
	if err != nil {
		return &pb.UpdateResp{
			Success: false,
		}, ErrUpdateUser
	}
	return &pb.UpdateResp{
		Success: true,
	}, nil
}
func (s *UserServiceService) GetUserInfo(ctx context.Context, req *pb.GetReq) (*pb.GetResp, error) {
	user, err := s.userHandler.GetUserInfoByID(ctx, req.GetUserId())
	if errors.Is(err, errcode.UserNotFound) {
		return &pb.GetResp{}, ErrUserNotFound
	}
	if err != nil {
		return &pb.GetResp{}, ErrGetUserInfo
	}
	return &pb.GetResp{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
		Age:   user.Age,
		Addr1: user.Addr1,
		Addr2: user.Addr2,
	}, nil
}
func (s *UserServiceService) SendVerifyCode(ctx context.Context, req *pb.SendReq) (*pb.SendResp, error) {
	if s.userHandler.CheckEmailExist(ctx, req.GetEmail()) {
		return &pb.SendResp{}, ErrEmailExist
	}
	code, err := s.userHandler.SendVerifyCode(ctx, req.GetEmail())
	if err != nil {
		return &pb.SendResp{}, ErrSendVerifyCode
	}
	return &pb.SendResp{
		Code: code,
	}, nil
}
