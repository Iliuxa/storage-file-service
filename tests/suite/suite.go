package suite

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"storage-file-service/internal/config"
	proto "storage-file-service/proto/gen"
	"strconv"
	"testing"
)

type Suite struct {
	*testing.T
	Cfg           *config.Config
	StorageClient proto.FileServiceClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()

	cfg := config.MustLoad("./../config/config.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)
	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	grpcAddress := net.JoinHostPort("localhost", strconv.Itoa(cfg.GRPC.Port))

	cc, err := grpc.DialContext(context.Background(),
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	storageClient := proto.NewFileServiceClient(cc)

	return ctx, &Suite{
		T:             t,
		Cfg:           cfg,
		StorageClient: storageClient,
	}
}
