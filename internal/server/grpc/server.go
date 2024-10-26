package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/VikaPaz/pantheon/internal/models"
	pb "github.com/VikaPaz/pantheon/proto/user"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type UserServise interface {
	CreateUser(ctx *context.Context, user models.User) (models.User, error)
	DeleteUser(ctx *context.Context, user models.User) error
	GetById(ctx *context.Context, user models.User) (models.User, error)
	GetByUsername(ctx *context.Context, user models.User) (models.User, error)
}

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	router   *grpc.Server
	service  UserServise
	validate *validator.Validate
	log      *logrus.Logger
}

func NewUserHandler(svc UserServise, log *logrus.Logger) *UserHandler {
	router := grpc.NewServer()
	validate := validator.New()
	return &UserHandler{
		router:   router,
		service:  svc,
		validate: validate,
		log:      log,
	}
}

func Run(server *UserHandler, port string) {
	server.router.RegisterService(&pb.UserService_ServiceDesc, server)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	err = server.router.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUsersResponse, error) {
	name := req.Username

	user := models.User{
		Username: name,
	}

	res, err := h.service.CreateUser(&ctx, user)
	if err != nil {
		return &pb.CreateUsersResponse{}, fmt.Errorf("failed to list Users: %w", err)
	}

	return &pb.CreateUsersResponse{
		Users: &pb.User{
			UserId:   res.Id,
			Username: req.Username,
			Coins:    res.Coins,
		},
	}, nil
}

func (h *UserHandler) GetUserById(ctx context.Context, req *pb.GetUsersByIdRequest) (*pb.GetUsersByIdResponse, error) {
	user, err := h.service.GetById(&ctx, models.User{Id: req.UserId})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &pb.GetUsersByIdResponse{
		Users: &pb.User{
			UserId:   user.Id,
			Username: user.Username,
			Coins:    user.Coins,
		},
	}, nil
}

func (h *UserHandler) GetUserByUsername(ctx context.Context, req *pb.GetUserByUsernameRequest) (*pb.GetUserByUsernameResponse, error) {
	user, err := h.service.GetByUsername(&ctx, models.User{Username: req.Username})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &pb.GetUserByUsernameResponse{
		Users: &pb.User{
			UserId:   user.Id,
			Username: user.Username,
			Coins:    user.Coins,
		},
	}, nil
}

func (h *UserHandler) Delete(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	user := models.User{Id: req.UserId}
	if err := h.service.DeleteUser(&ctx, user); err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	return &pb.DeleteUserResponse{UserId: user.Id}, nil
}
