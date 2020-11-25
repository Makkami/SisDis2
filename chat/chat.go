package chat

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
)

type Server struct {
	
}


func (s *Server) SubirChunk(ctx context.Context, message *Chunk) (*Message, error) {
	// write to disk
	fileName := message.Nombre + "_" + message.Parte
	_, err := os.Create(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// write/save buffer to disk
	ioutil.WriteFile(fileName, message.Buffer, os.ModeAppend)

	fmt.Println("Dividido en: ", fileName)
	return &Message{Body: ""}, nil
}