package main

import (
	"context"
	"fmt"
	pb "go-grpc-server/proto"
	"log"
	"net"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	PORT = 3001
	DB   = "grpc-comments"
	COLL = "comments"
)

type comment struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Comment string `json:"comment"`
}

type server struct {
	pb.UnimplementedGetInfoServer
}

func storeComment(comment string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mongoclient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/?readPreference=primary&appname=MongoDB%20Compass&ssl=false"))
	if err != nil {
		log.Fatal(err)
	}

	err = mongoclient.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	database := mongoclient.Database(DB)
	collection := database.Collection(COLL)

	var bdoc interface{}
	errb := bson.UnmarshalExtJSON([]byte(comment), true, &bdoc)
	fmt.Println(errb)

	insertResult, err := collection.InsertOne(ctx, bdoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertResult)
}

func (s *server) ReturnInfo(ctx context.Context, in *pb.RequestId) (*pb.ReplyInfo, error) {
	storeComment(in.GetId())
	fmt.Println(">> Received comment: %v", in.GetId)
	return &pb.ReplyInfo{Info: ">> I've received comment: " + in.GetId()}, nil
}

func main() {
	listen, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", PORT))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGetInfoServer(s, &server{})
	fmt.Printf("Server LIVE on http://localhost:%d", PORT)
	if err := s.Serve(listen); err != nil {
		log.Fatal("Server failed: %v", err)
	}

}
