package grpc

import (
	"fmt"
	"net"

	"github.com/nafisalfiani/p3-ugc-7-8/account-service/config"
	"github.com/nafisalfiani/p3-ugc-7-8/account-service/usecase"

	"github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPC interface {
	Run()
}

type grpcServer struct {
	cfg    *config.Value
	log    *logrus.Logger
	server *grpc.Server
}

func Init(cfg *config.Value, log *logrus.Logger, uc *usecase.Usecases) GRPC {
	s := grpc.NewServer()

	RegisterUserServiceServer(s, initUserGrpcServer(log, uc.User))

	return &grpcServer{
		cfg:    cfg,
		server: s,
		log:    log,
	}
}

func (g *grpcServer) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", g.cfg.Server.Base, g.cfg.Server.Port))
	if err != nil {
		g.log.Fatal(err)
	}

	g.log.Info("Listening and Serving GRPC on :50051")
	if err := g.server.Serve(listener); err != nil {
		g.log.Fatal(err)
	}
}

// wrapper to connect to grpc package
func Dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(target, opts...)
}

// wrapper to connect to grpc package
func WithInsecure() grpc.DialOption {
	return grpc.WithTransportCredentials(insecure.NewCredentials())
}
