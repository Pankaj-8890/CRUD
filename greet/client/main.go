package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "go-grpc/greet/proto"
	"log"
	"strconv"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserDetails struct {
	Id         int64 `json: "id"`
	First_name string `json: "first_name"`
	Second_name string `json: "second_name"`
	Age        int64	`json: "age"`
}

type CreateResponse struct {
	Token   string `json:"token"`
	Message string  `json:"message"`
}

type GetResponse struct{
	Message string `json:"message"`
}

type UpdateResponse struct{
	Message string `json:"message"`
}


var addr string = "localhost:50051"

func main(){
	conn,err := grpc.Dial(addr,grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil{
		log.Fatalf("Failed to connect %v\n",err)
	}

	defer conn.Close()


	c := pb.NewGreetClient(conn)

	router := mux.NewRouter()

	router.HandleFunc("/user",func(w http.ResponseWriter, r *http.Request){
		Create(c, w, r)}).Methods("POST")

	router.HandleFunc("/user/{id}",func(w http.ResponseWriter, r *http.Request){
		GetUser(c, w, r)}).Methods("GET")	

	router.HandleFunc("/user/{id}",func(w http.ResponseWriter, r *http.Request){
		UpdateUser(c, w, r)}).Methods("PUT")


	http.ListenAndServe(":8088",router)	

}


func GetUser(client pb.GreetClient, w http.ResponseWriter, r *http.Request){

	param := mux.Vars(r)

	userId, err := strconv.Atoi(param["id"])

	if err != nil {
		panic(err)
	}

	res,err := client.GetUser(context.Background(),&pb.GetUserRequest{
		Id: int64(userId),
	})

	if err!=nil{
		response := GetResponse{
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	usr := UserDetails{
		Id: res.User.Id,
		First_name: res.User.FirstName,
		Second_name: res.User.SecondName,
		Age: res.User.Age,
	}

	json.NewEncoder(w).Encode(usr)

}


func Create(client pb.GreetClient, w http.ResponseWriter, r *http.Request){

	
	var usr UserDetails

	json.NewDecoder(r.Body).Decode(&usr)

	fmt.Print(usr)

	user := &pb.User{
		Id : usr.Id,
		FirstName: usr.First_name,
		SecondName: usr.Second_name,
		Age: usr.Age,
	}
	res,err := client.CreatUser(context.Background(),&pb.CreateUserRequest{
		User: user,
	})
	
	if err!=nil{
		response := CreateResponse{
			Token:   "",
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := CreateResponse{
		Token:   res.Token,
		Message: res.Message,
	}
	json.NewEncoder(w).Encode(response)

}

func UpdateUser(client pb.GreetClient, w http.ResponseWriter, r *http.Request){

	
	param := mux.Vars(r)

	userId, err := strconv.Atoi(param["id"])

	if err != nil {
		panic(err)
	}

	var usr UserDetails

	json.NewDecoder(r.Body).Decode(&usr)

	user := &pb.User{
		Id : usr.Id,
		FirstName: usr.First_name,
		SecondName: usr.Second_name,
		Age: usr.Age,
	}
	res,err := client.UpdateUser(context.Background(),&pb.UpdateUserRequest{
		User: user,
		Id: int64(userId),
	})
	
	if err!=nil{
		response := UpdateResponse{
			Message: err.Error(),
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := UpdateResponse{
		Message: res.Message,
	}
	json.NewEncoder(w).Encode(response)

}


