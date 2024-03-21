package services

import (
	"fmt"

	"DOCKER_TEST/api-gateway/config"
	pbc "DOCKER_TEST/api-gateway/genproto/comment-service"
	pbp "DOCKER_TEST/api-gateway/genproto/post-service"
	pbu "DOCKER_TEST/api-gateway/genproto/user-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() pbu.UserServiceClient
	PostService() pbp.PostServiceClient
	CommentService() pbc.PostServiceClient
}

type serviceManager struct {
	userService    pbu.UserServiceClient
	postService    pbp.PostServiceClient
	commentService pbc.PostServiceClient
}

// user
func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}

// post
func (s *serviceManager) PostService() pbp.PostServiceClient {
	return s.postService
}

// comment
func (s *serviceManager) CommentService() pbc.PostServiceClient {
	return s.commentService
}

// user
func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")
	//coonuser
	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	//coonpost
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.PostServiceHost, conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	//cooncommnet
	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.CommenterviceHost, conf.CommentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		userService:    pbu.NewUserServiceClient(connUser),
		postService:    pbp.NewPostServiceClient(connPost),
		commentService: pbc.NewPostServiceClient(connComment),
	}

	return serviceManager, nil
}
