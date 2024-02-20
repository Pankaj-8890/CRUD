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
 
var DB *gorm.DB
var err error

type User struct {
	gorm.Model
	Id         int64
	First_name string
	Second_name string
	Age        int64
	Token      string
}

func DatabaseConnection() {
	host := "localhost"
	port := "5432"
	dbName := "postgres"
	dbUser := "postgres"
	password := "postgres"
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",host,port,dbUser,dbName,password)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(&User{})

	if err != nil {
		log.Fatal("Error connecting to the database...", err)
	}
	fmt.Println("Database connection successful...")

 }

var addr string = "0.0.0.0:50051"

type Server struct{
	pb.GreetServer
}

func main(){
	lis, err := net.Listen("tcp",addr)

	if err != nil{
		log.Fatalf("Failed to Listen :%v\n",err)
	}

	log.Printf("listening %s\n",addr)

	s := grpc.NewServer()

	pb.RegisterGreetServer(s,&Server{})

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

	// fmt.Println("-----------------")
	// fmt.Println(users)

	res := DB.Create(&users)
	if res.RowsAffected == 0 {
		return nil, errors.New("movie creation unsuccessful")
	}

	response := &pb.CreateUserResponse{
		Token:   users.Token,
		Message: "User successfully created",
	}

	return response,nil

}

// func (s *Server)GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error){


// }

// func (s *Server)UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error){

// }