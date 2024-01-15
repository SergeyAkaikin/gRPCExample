package main

import (
	"context"
	pb "gRPCEx/gen/go"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"
)

const address = "grpc:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("didnt connect %v", err)
	}

	defer conn.Close()

	c := pb.NewProductInfoClient(conn)

	name := "Apple iphone 11"
	description := "some text"
	price := float32(1000.0)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatalf("Couldn't add product: %v", err)
	}

	log.Printf("Product ID: %s added successfully", r.Value)

	name = "Apple ipad 11"
	description = "some text2"
	price = float32(10000.0)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err = c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatalf("Couldn't add product: %v", err)
	}

	updateID1 := r.Value
	log.Printf("Product ID: %s added successfully", r.Value)

	name = "Samsung galaxy note"
	description = "some text3"
	price = float32(1000.0)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err = c.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil {
		log.Fatalf("Couldn't add product: %v", err)
	}

	updateID2 := r.Value
	log.Printf("Product ID: %s added successfully", r.Value)

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	product, err := c.GetProduct(ctx, &pb.ProductID{Value: r.Value})
	if err != nil {
		log.Fatalf("Could not get product: %v", err)
	}

	log.Printf("Product: %s", product.String())

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	searchStream, _ := c.SearchProduct(ctx, &wrappers.StringValue{Value: "Apple"})

	for {
		searchOrder, err := searchStream.Recv()
		if err == io.EOF {
			break
		}

		log.Println("Search result: ", searchOrder)
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	updateStream, err := c.UpdateProduct(ctx)

	if err != nil {
		log.Fatalf("%v.UpdateProduct() = _, %v", c, err)
	}

	updateP1 := pb.Product{Id: updateID1, Name: "Apple updated ipad 11", Description: "updated desc"}
	updateP2 := pb.Product{Id: updateID2, Name: "Samsung updated galaxy note", Description: "updated desc"}
	if err := updateStream.Send(&updateP1); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updateP1, err)
	}

	if err := updateStream.Send(&updateP2); err != nil {
		log.Fatalf("%v.Send(%v) = %v", updateStream, updateP1, err)
	}

	updateRes, err := updateStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() error %v, want %v", updateStream, err, nil)
	}

	log.Println("Update products:", updateRes)
}
