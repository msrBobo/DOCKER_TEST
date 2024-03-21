package postgres

// import (
// 	"DOCKER_TEST/user-service/config"
// 	pb "DOCKER_TEST/user-service/genproto/user-service"
// 	"DOCKER_TEST/user-service/pkg/db"
// 	"DOCKER_TEST/user-service/storage/repo"
// 	"log"
// 	"testing"

// 	"github.com/google/uuid"

// 	"github.com/stretchr/testify/suite"
// )

// type UserReposisitoryTestSuite struct {
// 	suite.Suite
// 	CleanUpFunc func()
// 	Repository  repo.UserStorageI
// }

// func (s *UserReposisitoryTestSuite) SetupSuite() {
// 	pgPoll, err, cleanUp := db.ConnectToDB(config.Load())
// 	if err != nil {
// 		log.Fatal("Error while connecting database with suite test")
// 		return
// 	}
// 	s.CleanUpFunc = cleanUp
// 	s.Repository = NewUserRepo(pgPoll)
// }

// // test func
// func (s *UserReposisitoryTestSuite) TestUserCRUD() {
// 	// struct for create user
// 	user := &pb.User{
// 		Name:     "test name",
// 		LastName: "test last name",
// 	}

// 	// uuid generating
// 	user.Id = uuid.New().String()

// 	// check create user method
// 	createdUser, err := s.Repository.Create(user)
// 	s.Suite.NotNil(createdUser)
// 	s.Suite.NoError(err)
// 	s.Suite.Equal(user, createdUser)

// 	// struct for get user method
// 	userRequst := &pb.UserRequest{
// 		UserId: createdUser.Id,
// 	}
// 	// check get user method
// 	getCreatedUser, err := s.Repository.Get(userRequst)
// 	s.Suite.NoError(err)
// 	s.Suite.NotNil(getCreatedUser)
// 	s.Suite.Equal(getCreatedUser, user)
// 	s.Suite.Equal(getCreatedUser, createdUser)

// 	// check update user method
// 	updatedUser, err := s.Repository.Update(userRequst)
// 	s.Suite.NoError(err)
// 	s.Suite.NotNil(updatedUser)
// 	s.Suite.NotEqual(updatedUser, getCreatedUser)

// 	// check get all users method
// 	getAllRequest := &pb.GetAllUsersRequest{
// 		Page:  1,
// 		Limit: 30,
// 	}
// 	getAllUsers, err := s.Repository.GetAll(getAllRequest)
// 	s.Suite.NoError(err)
// 	s.Suite.NotNil(getAllUsers)

// 	// check delete user method
// 	deletedUser, err := s.Repository.Delete(userRequst)
// 	s.Suite.NoError(err)
// 	s.Suite.NotNil(deletedUser)
// }

// func TestExampleTestSuite(t *testing.T) {
// 	suite.Run(t, new(UserReposisitoryTestSuite))
// }
