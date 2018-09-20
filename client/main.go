package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	naruto "../services/proto/rpc/src"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	address = "localhost:8000"
)

// Metadata Metadata of grpc headers
type Metadata map[string][]string
type grpcClient struct{}

/**
 * Add functions here. All parameters should be string and
 * converted to required type, in order that the functions can
 * be called via reflection.
 */

// GetRoom Fetch room information from roomsvr
func (c grpcClient) GetRoom(roomID string) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error creating connection.\nErr: %v", err)
	}
	defer conn.Close()

	stub := naruto.NewRoomServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "server-family", "roomsvr")

	roomInt, err := strconv.Atoi(roomID)
	response, err := stub.GetRoom(ctx, &naruto.RoomRequest{Id: int32(roomInt)})
	if err != nil {
		fmt.Printf("Error sending request.\nErr: %v", err)
	}
	fmt.Printf("Response: %s", response)
}

// ReflectValues Generate value array of reflect.Value from string array
func ReflectValues(arr []string) []reflect.Value {
	result := make([]reflect.Value, len(arr))
	for i := 0; i < len(arr); i++ {
		result[i] = reflect.ValueOf(arr[i])
	}
	return result
}

const (
	prompt = `Missing parameters.
Usage:
	grpc-client [funcName] [...args]
	grpc-client -m
	`
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf(prompt)
		return
	}
	funcName := os.Args[1]
	args := os.Args[2:]
	grpc := grpcClient{}
	grpcValue := reflect.ValueOf(grpc)
	grpcValue.MethodByName(funcName).Call(ReflectValues(args))
}
