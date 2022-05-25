package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "go-grpc-traffic-client/proto"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

const (
	addr = "localhost:3001"
)

func connectSv(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"mensaje\":\"OK\"}"))
		return
	}

	data, _ := ioutil.ReadAll(r.Body)

	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		json.NewEncoder(w).Encode("Error connecting to server.")
		log.Fatalf("Error connecting to server. (%v)", err)
	}

	defer conn.Close()

	cl := pb.NewGetInfoClient(conn)

	id := string(data)
	if len(os.Args) > 1 {
		id = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ret, err := cl.ReturnInfo(ctx, &pb.RequestId{Id: id})
	if err != nil {
		json.NewEncoder(w).Encode("Couldn't retrieve information.")
		log.Fatalf("Couldn't retrieve information. (%v)", err)
	}

	log.Printf("SV response: %s\n", ret.GetInfo())
	json.NewEncoder(w).Encode("Se ha almacenado la informacion")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", connectSv)
	fmt.Println("Client LIVE on http://localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", router))
}
