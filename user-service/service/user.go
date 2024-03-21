package service

import (
	pb "DOCKER_TEST/user-service/genproto/user-service"
	l "DOCKER_TEST/user-service/pkg/logger"
	"DOCKER_TEST/user-service/storage"
	"context"

	"github.com/jmoiron/sqlx"
)

// UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	pb.UnimplementedUserServiceServer
}

// NewUserService ...
func NewUserService(db *sqlx.DB, log l.Logger) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

func (s *UserService) Create(ctx context.Context, req *pb.UserCreateReq) (*pb.UserResponse, error) {
	user, err := s.storage.User().Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Get(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := s.storage.User().Get(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *pb.Email) (*pb.UserResponse, error) {
	user, err := s.storage.User().GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetAll(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := s.storage.User().GetAll(req)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserService) Update(ctx context.Context, user *pb.UserUpd) (*pb.UserResponse, error) {
	resp, err := s.storage.User().Update(user)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *UserService) Delete(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	user, err := s.storage.User().Delete(req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CheckUniquess(ctx context.Context, req *pb.CheckUniqReq) (*pb.CheckUniqResp, error) {
	// Call the data access layer to check uniqueness based on email.
	// Assuming you have a context variable named "ctx" and a request variable named "req"
	exists, err := s.storage.User().CheckUniquenessByEmail(req)

	if err != nil {
		// Handle error
		return nil, err
	}
	return &pb.CheckUniqResp{
		IsExist: exists,
	}, nil

}
