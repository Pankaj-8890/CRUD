package main

import (
	"context"
	// "errors"
	"testing"

	pb "go-grpc/greet/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MockGreetServer struct {
    DB *gorm.DB
	mock.Mock
}


func (m *MockGreetServer) CreatUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.CreateUserResponse), args.Error(1)
}

func TestCreatUser(t *testing.T) {


	mockDb, mock, _ := sqlmock.New()
    dialector := postgres.New(postgres.Config{
    Conn:       mockDb,
    DriverName: "postgres",
    })
    db, _ := gorm.Open(dialector, &gorm.Config{})

    rows := sqlmock.NewRows([]string{"Code", "Price"}).AddRow("D43", 100)
    mock.ExpectQuery(`SELECT`).WillReturnRows(rows)

   
    req := &pb.CreateUserRequest{
        User: &pb.User{
            Id:         1,
            FirstName:  "John",
            SecondName: "Doe",
            Age:        30,
        },
    }

    mockServer := &MockGreetServer{
        DB : db,
    }

	mockServer.On("CreatUser", context.Background(), req).Return(&pb.CreateUserResponse{
		Message: "User successfully created",
	}, nil)

	res, err := mockServer.CreatUser(context.Background(), req)

	assert.NoError(t, err, "Unexpected error in CreatUser")
	assert.NotNil(t, res, "Response should not be nil")
	assert.Equal(t, "User successfully created", res.Message, "Unexpected response message")

	mockServer.AssertExpectations(t)

}
