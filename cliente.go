package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"log"
	"net"
	"math"
	"bufio"
	"strings"
	"math/rand"
	"time"

	"github.com/Makkami/SisDis2/chat"
	"google.golang.org/grpc"
)	


func crearServer() {
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

	go crearServer()

	var request = bufio.NewReader(os.Stdin)
    for {

		fmt.Println("----------------------------")
        fmt.Printf("Ingrese una opcion de Orden: ")
        libro, _ := request.ReadString('\n')
        libro = strings.Trim(libro, " \r\n")
        fileToBeChunked := "./" + libro + ".pdf"

		flag1, flag2, flag3 := true, true, true
		/* Conexiones a los datanodes*/
			// Datanode 1
		var conn *grpc.ClientConn
		conn, err := grpc.Dial("dist137:9001", grpc.WithInsecure())
		if err != nil {
			flag1 = false
			fmt.Printf("No se pudo conectar al DataNode 1:  %s", err)
		}
		c := chat.NewChatServiceClient(conn)

		defer conn.Close()

			// Datanode 2
		var conn2 *grpc.ClientConn
		conn2, err2 := grpc.Dial("dist138:9002", grpc.WithInsecure())
		flag2 = false
		if err2 != nil {
			flag2 = false
			fmt.Printf("No se pudo conectar al DataNode 2:  %s", err2)
		}
		c2 := chat.NewChatServiceClient(conn2)

		defer conn2.Close()

			// Datanode 3
		var conn3 *grpc.ClientConn
		conn3, err3 := grpc.Dial("dist139:9003", grpc.WithInsecure())
		if err3 != nil {
			flag3 = false
			fmt.Printf("No se pudo conectar al DataNode 3:  %s", err3)
		}
		c3 := chat.NewChatServiceClient(conn3)

		defer conn3.Close()
	

		//Elegir Datanode random
		sliceDN := []int{1,2,3}
		if (!flag1 && !flag2 && !flag3) {
			fmt.Println("No se pudo establecer conexión que ningun DataNode")
			os.Exit(1)
		} else if (!flag1 && !flag2) {sliceDN = []int{3}
		} else if !flag1 && !flag3 {sliceDN = []int{2}
		} else if !flag2 && !flag3 {sliceDN = []int{1}
		} else if !flag1 {sliceDN = []int{2,3}
		} else if !flag2 {sliceDN = []int{1,3}
		} else if !flag3 {sliceDN = []int{1,2}}
		
		rand.Seed(time.Now().UnixNano())
		rdIndex := rand.Intn(len(sliceDN))
		dn_rand := sliceDN[rdIndex]
		fmt.Printf("Random: %d\n", dn_rand)

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
				Nombre: libro,
				Parte: strconv.FormatUint(i, 10),
				NumPartes: totalPartsNum,
				Buffer: partBuffer,
			}

			var response *chat.Message
			switch dn_rand {
			case 1:
				response, _ = c.SubirChunk(context.Background(), &message)
				log.Printf("DataNode1 %s", response.Body)
			case 2:
				response, _ = c2.SubirChunk(context.Background(), &message)
				log.Printf("DataNode2 %s", response.Body)
			case 3:
				response, _ = c3.SubirChunk(context.Background(), &message)
				log.Printf("DataNode3 %s", response.Body)
			}
		}
	}
}