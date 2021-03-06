// Package chunks implements a client for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	pb "lab2"

	"google.golang.org/grpc"
)

// constantes de puertos y nombres de instancias
const (
	address    = "dist29:50051"
	clientName = "Sender"
)

func uploadBook(fileName string, c pb.GreeterClient) {

	fileToBeChunked := "./" + fileName

	file, err := os.Open(fileToBeChunked)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 250000 // 1 * (1 << 20) // 1 MB, change this to your requirement

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		partBuffer := make([]byte, partSize)

		file.Read(partBuffer)

		// generacion de orden
		storeRequest := &pb.StoreRequest{
			ChunkPart:  strconv.FormatUint(i, 10),
			FileName:   fileName,
			ClientName: clientName,
			Chunk:      partBuffer,
			Part:       int32(totalPartsNum),
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// Hacer una consulta
		r, err := c.ReceiveChunk(ctx, storeRequest)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("%s", r)
	}
}

func main() {

	// Set up a connection to the server.
	// Contact the server and print out its response.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	uploadBook("Don_Quijote_de_la_Mancha-Cervantes_Miguel.pdf", c)
}
