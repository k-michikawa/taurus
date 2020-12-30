package infrastructures

import (
	"log"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// フレームワークに従ってgRPCサーバーを生成する
func CreateServer() *grpc.Server {
	contextMiddleware := generateContextMiddleware()
	loggingMiddleware := generateLoggingMiddleware()
	recoveryMiddleware := generateRecoveryMiddleware()

	return grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			contextMiddleware,
			loggingMiddleware,
		)),
		grpc_middleware.WithUnaryServerChain(grpc_middleware.ChainUnaryServer(
			recoveryMiddleware,
		)),
	)
}

// リクエスト情報をMetadataに埋め込んでくれるミドルウェア生成
func generateContextMiddleware() grpc.StreamServerInterceptor {
	return grpc_ctxtags.StreamServerInterceptor()
}

// いい感じにアクセスログ出してくれるミドルウェア生成
func generateLoggingMiddleware() grpc.StreamServerInterceptor {
	zapOpts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}
	zapLogger, _ := zap.NewProduction()
	grpc_zap.ReplaceGrpcLogger(zapLogger)

	return grpc_zap.StreamServerInterceptor(zapLogger, zapOpts...)
}

// パニック起こしたときの処理を登録しておくミドルウェア生成
func generateRecoveryMiddleware() grpc.UnaryServerInterceptor {
	recoveryFunc := func(p interface{}) error {
		// とりあえずログ吐いて500相当のものを返して復帰
		log.Printf("Paniced!!!!!!!!!!\n%+v\n", p)
		return grpc.Errorf(codes.Internal, "Unexpected error")
	}
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(recoveryFunc),
	}

	return grpc_recovery.UnaryServerInterceptor(recoveryOpts...)
}
