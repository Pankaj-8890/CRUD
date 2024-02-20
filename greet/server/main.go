package main

import (
	pb "go-grpc/greet/proto"
	"log"
	"net"
	"fmt"
	"google.golang.org/grpc"
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"errors"
)

func init() {
	DatabaseConnection()
 }

type Server struct{
	DB *gorm.DB
	pb.GreetServer
}
 
// var DB *gorm.DB
var err error

type User struct {
	gorm.Model
	Id         int64
	First_name string
	Second_name string
	Age        int64
	Token      string
}



 func DatabaseConnection() *gorm.DB {
    host := "localhost"
    port := "5432"
    dbName := "postgres"
    dbUser := "postgres"
    password := "postgres"
    dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, dbUser, dbName, password)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Error connecting to the database...", err)
    }
    fmt.Println("Database connection successful...")

    // Auto-migrate models
    db.AutoMigrate(&User{})

    return db
}


var addr string = "0.0.0.0:50051"



func main(){
	lis, err := net.Listen("tcp",addr)

	if err != nil{
		log.Fatalf("Failed to Listen :%v\n",err)
	}

	log.Printf("listening %s\n",addr)

	s := grpc.NewServer()
	db := DatabaseConnection()
	pb.RegisterGreetServer(s,&Server{DB : db})

	if err = s.Serve(lis); err != nil{
		log.Fatalf("failed to server %v\n",err)
	}
	
}


	
	
func (s *Server)CreatUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error){
	

	usr := req.User
	token := uuid.New().String()

	
	users := User{
		Id: usr.Id,
		First_name: usr.FirstName,
		Second_name: usr.SecondName,
		Age: usr.Age,
	}

	users.Token = token 

	res := s.DB.Create(&users)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie creation unsuccessful")
	}

	response := &pb.CreateUserResponse{
		Token:   users.Token,
		Message: "User successfully created",
	}

	return response,nil

}

func (s *Server)GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error){

	id := req.Id

	var user *pb.User
	s.DB.First(&user, id)
	
	if user.Id == 0 {
		return nil, errors.New("user not found")
	}

	response := &pb.GetUserResponse{
		User: user,
	}
	return response,nil

}

// func (s *Server)UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error){

// }

// mockgen -destination=mocks/mock_server.go -package=mocks go-grpc/greet/server/main.go
// mockgen -source=greet/server/main.go -destination=mocks/mock_server.go -package=mocks github.com/jinzhu/gorm DB
