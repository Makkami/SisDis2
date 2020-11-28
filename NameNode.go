package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	"strings"

	"github.com/Makkami/SisDis2/chat"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error al conectar con el DataNode1: %s", err)
	}
	defer conn.Close()

	chat := chat.NewChatServiceClient(conn)
}