package main

import (
	"context"
	"fmt"
	pb "gRPCEx/gen/go"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const port = ":50051"

type server struct {
	pb.UnimplementedProductInfoServer
	productMap map[string]*pb.Product
	r          *rand.Rand
}

func (s *server) AddProduct(ctx context.Context, in *pb.Product) (*pb.ProductID, error) {
	if s.r == nil {
		s.r = rand.New(rand.NewSource(time.Now().Unix()))
	}

	out := strconv.Itoa(s.r.Int())

	in.Id = out

	if s.productMap == nil {
		s.productMap = make(map[string]*pb.Product)
	}
	s.productMap[out] = in

	return &pb.ProductID{Value: out}, status.New(codes.OK, "").Err()
}
func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.Product, error) {
	value, exists := s.productMap[in.Value]
	if exists {
		return value, status.New(codes.OK, "").Err()
	}

	return nil, status.Errorf(codes.NotFound, "Product does not exists", in.Value)
}

func (s *server) SearchProduct(searchQuery *wrappers.StringValue, stream pb.ProductInfo_SearchProductServer) error {
	for key, product := range s.productMap {
		log.Print(key, product)
		if strings.Contains(product.Name, searchQuery.Value) {
			if err := stream.Send(product); err != nil {
				return fmt.Errorf("error sending message to stream: %v")
			}
			log.Println("Matching product found: " + product.Name)
		}
	}
	return nil
}

func (s *server) UpdateProduct(stream pb.ProductInfo_UpdateProductServer) error {
	productsStr := "Updated products IDs: "

	for {
		product, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&wrappers.StringValue{Value: "Products processed " + productsStr})
		}
		s.productMap[product.Id] = product

		log.Println("Product ID", product.Id, ": Updated")
		productsStr += product.Id + ", "
	}
}
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen")
	}
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	log.Printf("Starting gRPC listener on port" + port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
