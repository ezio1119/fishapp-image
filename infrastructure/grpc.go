package infrastructure

import (
	"github.com/ezio1119/fishapp-image/infrastructure/middleware"
	"github.com/ezio1119/fishapp-image/pb"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGrpcServer(middL middleware.Middleware, imageController pb.ImageServiceServer) *grpc.Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middL.UnaryLogingInterceptor(),
			// middL.UnaryRecoveryInterceptor(),
			middL.UnaryValidationInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			middL.StreamLogingInterceptor(),
			// middL.StreamRecoveryInterceptor(),
			middL.StreamValidationInterceptor(),
		)),
	)

	pb.RegisterImageServiceServer(server, imageController)
	reflection.Register(server)
	return server
}
