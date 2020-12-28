package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "taurus/pb/user"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// Server is gRPC server
type Server struct {
	pb.UserServiceServer
}

// PostUser Service
func (s *Server) PostUser(c context.Context, r *pb.PostUserRequest) (*pb.PostUserResponse, error) {
	log.Printf("%v", c)
	return &pb.PostUserResponse{
		User: &pb.User{
			Id:        "id",
			Name:      "name",
			Email:     "email",
			CreatedAt: 0,
			UpdatedAtOneof: &pb.User_UpdatedAt{
				UpdatedAt: 0,
			},
		},
	}, nil
}

// ListUser Service
func (s *Server) ListUser(c context.Context, r *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	log.Printf("%v", c)
	return &pb.ListUserResponse{
		Users: []*pb.User{},
	}, nil
}

// ReadUser Service
func (s *Server) ReadUser(c context.Context, r *pb.ReadUserRequest) (*pb.ReadUserResponse, error) {
	log.Printf("%v", c)
	return &pb.ReadUserResponse{
		User: &pb.User{
			Id:        "id",
			Name:      "name",
			Email:     "email",
			CreatedAt: 0,
			UpdatedAtOneof: &pb.User_UpdatedAt{
				UpdatedAt: 0,
			},
		},
	}, nil
}

// UpdateUser Service
func (s *Server) UpdateUser(c context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Printf("%v", c)
	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:        "id",
			Name:      "name",
			Email:     "email",
			CreatedAt: 0,
			UpdatedAtOneof: &pb.User_UpdatedAt{
				UpdatedAt: 0,
			},
		},
	}, nil
}

// DeleteUser Service
func (s *Server) DeleteUser(c context.Context, r *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Printf("%v", c)
	return &pb.DeleteUserResponse{}, nil
}

// server running
func run(port string) error {
	log.Print("Init server...")
	// listenポートの取得
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "Failed listen port")
	}

	// ログ関連のオプション
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}
	zapLogger, _ := zap.NewProduction()
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	// リカバリー関連のオプション
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	// Serverを作ってServiceを登録
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			// リクエスト情報をMetadataに埋め込んでくれるミドルウェア
			grpc_ctxtags.StreamServerInterceptor(),
			// いい感じにアクセスログ出してくれるミドルウェア
			grpc_zap.StreamServerInterceptor(zapLogger, zapOpts...),
		)),
		grpc_middleware.WithUnaryServerChain(grpc_middleware.ChainUnaryServer(
			// パニック起こしたときにrecoveryFunc読んでくれるようにしてくれるミドルウェア
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		)),
	)
	var server Server
	pb.RegisterUserServiceServer(s, &server)

	// サーバーを起動
	log.Printf("Start taurus Server port: %s", port)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "Failed running server")
	}

	return nil
}

// パニック起こしたときにgRPCのコードに変換してレスポンスを返す
func recoveryFunc(p interface{}) error {
	log.Printf("p: %+v\n", p)
	return grpc.Errorf(codes.Internal, "Unexpected error")
}

func main() {
	if err := run(":9010"); err != nil {
		log.Fatalf("%v", err)
	}
}
