package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"log"
	"net"
	"math"

	"github.com/Makkami/SisDis2/chat"
	"google.golang.org/grpc"
)	


func con() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := chat.Server{}
	grpcServer := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServer, &s)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("a %v", err)
	}
}



func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	defer conn.Close()
	c := chat.NewChatServiceClient(conn)


	fileToBeChunked := "./Mujercitas-Alcott_Louisa_May.pdf"

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 250 * (1 << 10) // Este (1 << 10) es igual a 2^10, entonces es 250 * 1024 = 256000

	// calculate total number of parts the file will be chunked into

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Dividiendo el archivo en %d partes.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		message := chat.Chunk{
			Nombre: "Mujercitas",
			Parte: strconv.FormatUint(i, 10),
			NumPartes: totalPartsNum,
			Buffer: partBuffer,
		}

		var response *chat.Message
		
		response, _ = c.SubirChunk(context.Background(), &message)
		log.Printf("Aca %s", response.Body)
	}
}