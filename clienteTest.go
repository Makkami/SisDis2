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
		/* Conexiones a los datanodes*/
			// Datanode 1
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(":9001", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("uwu %s", err)
		}
		c := chat.NewChatServiceClient(conn)

		defer conn.Close()

			// Datanode 2
		var conn2 *grpc.ClientConn
		conn2, err2 := grpc.Dial(":9002", grpc.WithInsecure())
		if err2 != nil {
			log.Fatalf("uwu %s", err2)
		}
		c2 := chat.NewChatServiceClient(conn2)

		defer conn2.Close()

			// Datanode 3
		var conn3 *grpc.ClientConn
		conn3, err3 := grpc.Dial(":9003", grpc.WithInsecure())
		if err3 != nil {
			log.Fatalf("uwu %s", err3)
		}
		c3 := chat.NewChatServiceClient(conn3)

		defer conn3.Close()

        fmt.Println("----------------------------")
		fmt.Printf("Ingrese una opcion de Orden\n")
		fmt.Printf("1. Subir libro\n")
		fmt.Printf("2. Descargar libro\n")
		opcion := request.ReadString('\n')
		if opcion == '1' {
			
			libro, _ := request.ReadString('\n')
			libro = strings.Trim(libro, " \r\n")
			fileToBeChunked := "./" + libro + ".pdf"
		
			/*
			//Elegir Datanode random
			rand.Seed(time.Now().UnixNano())
			dn_rand := rand.Intn(3) + 1
			fmt.Printf("Random: %d", dn_rand)
			*/

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

			reparto := uint64(totalPartsNum/3)

			for i := uint64(0); i < reparto; i++ {

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
				response, _ = c.SubirChunk(context.Background(), &message)
				log.Printf("DataNode1 %s", response.Body)
			}

			for i := uint64(0); i < reparto; i++ {

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
				response, _ = c2.SubirChunk(context.Background(), &message)
				log.Printf("DataNode2 %s", response.Body)
			}


			for i := uint64(0); i < reparto; i++ {

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
				response, _ = c3.SubirChunk(context.Background(), &message)
				log.Printf("DataNode3 %s", response.Body)
			}
		}


		if opcion == 2 {
			var writePosition int64 = 0
			for j := uint64(0); j < totalPartsNum; j++ {

				//read a chunk
				currentChunkFileName := "Mujercitas_" + strconv.FormatUint(j, 10)

				newFileChunk, err := os.Open(currentChunkFileName)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				defer newFileChunk.Close()

				chunkInfo, err := newFileChunk.Stat()

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// calculate the bytes size of each chunk
				// we are not going to rely on previous data and constant

				var chunkSize int64 = chunkInfo.Size()
				chunkBufferBytes := make([]byte, chunkSize)

				fmt.Println("Appending at position : [", writePosition, "] bytes")
				writePosition = writePosition + chunkSize

				// read into chunkBufferBytes
				reader := bufio.NewReader(newFileChunk)
				_, err = reader.Read(chunkBufferBytes)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				// DON't USE ioutil.WriteFile -- it will overwrite the previous bytes!
				// write/save buffer to disk
				//ioutil.WriteFile(newFileName, chunkBufferBytes, os.ModeAppend)

				n, err := file.Write(chunkBufferBytes)

				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				file.Sync() //flush to disk

				// free up the buffer for next cycle
				// should not be a problem if the chunk size is small, but
				// can be resource hogging if the chunk size is huge.
				// also a good practice to clean up your own plate after eating

				chunkBufferBytes = nil // reset or empty our buffer

				fmt.Println("Written ", n, " bytes")

				fmt.Println("Recombining part [", j, "] into : ", newFileName)
			}

			// now, we close the newFileName
			file.Close()
		}	
	}
}