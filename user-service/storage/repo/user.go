package repo

import (
	pb "DOCKER_TEST/user-service/genproto/user-service"
	"context"
)

// UserStorageI ...
type UserStorageI interface {
	Create(ctx context.Context, user *pb.UserCreateReq) (*pb.UserResponse, error)
	Get(request *pb.UserRequest) (*pb.UserResponse, error)
	GetAll(request *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error)
	Update(user *pb.UserUpd) (*pb.UserResponse, error)
	Delete(request *pb.UserRequest) (*pb.UserResponse, error)
	CheckUniquenessByEmail(req *pb.CheckUniqReq) (bool, error)
	GetUserByEmail(email string) (*pb.UserResponse, error)
}
