package main

import (
	"context"
	pb "go-grpc/greet/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
)

var lis *bufconn.Listener

func serverMock() {

	lis = bufconn.Listen(1024 * 1024)

	srv := grpc.NewServer()

	svc := Server{}
	pb.RegisterGreetServer(srv, &svc)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to server %v\n", err)
		}
	}()

}
func TestSayHello(t *testing.T) {

	serverMock()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	t.Cleanup(func() {
		conn.Close()
	})

	if err != nil {
		t.Fatalf("grpc.DialContext %v", err)
	}

	tests := []struct {
		name string
		want string
	}{
		{
			name: "world",
			want: "Helloworld",
		},
		
	}

	for _, tt := range tests {
		client := pb.NewGreetClient(conn)
		res, err := client.Greet(context.Background(), &pb.Request{FirstName: tt.name})
		if err != nil {
			// log.Fatalf("error 1 %v",err)
			t.Errorf("HelloTest(%v) got unexpected error",err)
		}
		if res.Result != tt.want {
			// log.Fatalf("error 2 %v",err)
			t.Errorf("HelloText(%v)=%v, wanted %v", tt.name, res.Result, tt.want)
		}
	}
}

