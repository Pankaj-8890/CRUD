package main

import (
    "context"
    // "errors"
    "testing"

    pb "go-grpc/greet/proto"
	"gorm.io/gorm"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
	// "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"github.com/golang/mock/gomock"
	
)


func TestCreatUser(t *testing.T) {


	ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockDB := mocks.NewMockDB(ctrl)

    mockDB.EXPECT().Open(gomock.Any(), gomock.Any()).Return(nil)

    // Create a test user request
    testUserRequest := &pb.CreateUserRequest{
        User: &pb.User{
            Id:         1,
            FirstName:  "John",
            SecondName: "Doe",
            Age:        30,
        },
    }

    // Create a test context
    ctx := context.Background()

    // Mock the behavior of DB.Create method
    mockDB.On("Create", mock.AnythingOfType("*main.User")).Return(&gorm.DB{
        RowsAffected: 1,
    })

    // Initialize the server with the mocked DB
    server := &Server{
        DB: &mockDB{},
    }

    // Call the function to be tested
    response, err := server.CreatUser(ctx, testUserRequest)

    // Assert the response and error
    assert.NoError(t, err)
    assert.NotNil(t, response)
    assert.Equal(t, "User successfully created", response.Message)
    assert.NotEmpty(t, response.Token)

    // Verify that the mock was called with the expected argument
    mockDB.AssertCalled(t, "Create", mock.AnythingOfType("*main.User"))
}

func TestCreatUser_Unsuccessful(t *testing.T) {
    // Initialize MockDB
    mockDB := new(MockDB)

    // Create a test user request
    testUserRequest := &pb.CreateUserRequest{
        User: &pb.User{
            Id:         1,
            FirstName:  "John",
            SecondName: "Doe",
            Age:        30,
        },
    }

    // Create a test context
    ctx := context.Background()

    // Mock the behavior of DB.Create method
    mockDB.On("Create", mock.AnythingOfType("*main.User")).Return(&gorm.DB{
        RowsAffected: 0,
    })

    // Initialize the server with the mocked DB
    server := &Server{
        DB: &mockDB{},
    }

    // Call the function to be tested
    response, err := server.CreatUser(ctx, testUserRequest)

    // Assert the error
    assert.Error(t, err)
    assert.Nil(t, response)

    // Verify that the mock was called with the expected argument
    mockDB.AssertCalled(t, "Create", mock.AnythingOfType("*main.User"))
}