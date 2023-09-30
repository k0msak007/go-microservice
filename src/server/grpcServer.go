package server

import (
	"log"
	"net"

	"github.com/k0msak007/go-microservice/src/config"
	grpccontrollers "github.com/k0msak007/go-microservice/src/controllers/grpcControllers"
	pbItem "github.com/k0msak007/go-microservice/src/proto/item"
	"github.com/k0msak007/go-microservice/src/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

type grpcServer struct {
	cfg      *config.Config
	dbClient *mongo.Client
}

func NewServerGrpc(cfg *config.Config, dbClient *mongo.Client) *grpcServer {
	return &grpcServer{
		cfg:      cfg,
		dbClient: dbClient,
	}
}

func (s *grpcServer) StartItemServer() {
	opts := make([]grpc.ServerOption, 0)
	gs := grpc.NewServer(opts...)

	log.Printf("%s server is starting on %s", s.cfg.App.Appname, s.cfg.App.Url)
	lis, err := net.Listen("tcp", s.cfg.App.Url)
	if err != nil {
		log.Fatal(err)
	}

	pbItem.RegisterItemServiceServer(
		gs,
		&grpccontrollers.ItemGrpcController{
			ItemRepository: repositories.ItemRepository{
				Client: s.dbClient,
			},
		},
	)
	gs.Serve(lis)
}
