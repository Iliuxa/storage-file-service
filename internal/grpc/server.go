package storagegrpc

import (
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"storage-file-service/internal/domain"
	"storage-file-service/internal/usecase"
	proto "storage-file-service/proto/gen"
)

type serverApi struct {
	proto.UnimplementedFileServiceServer
	storage usecase.StorageUsecase
}

func Register(gRPCServer *grpc.Server, storage usecase.StorageUsecase) {
	proto.RegisterFileServiceServer(gRPCServer, &serverApi{storage: storage})
}

func (s *serverApi) GetFileList(
	ctx context.Context,
	in *proto.FileListRequest,
) (*proto.FileListResponse, error) {
	files, err := s.storage.GetFileList(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrUsersGettingError) {
			return nil, status.Error(codes.InvalidArgument, "Failed to getting users")
		}
		return nil, status.Error(codes.Internal, "Failed to getting users")
	}
	return &proto.FileListResponse{Files: files}, nil
}

func (s *serverApi) UploadFile(
	str grpc.ClientStreamingServer[proto.UploadRequest, proto.UploadResponse],
) error {

	return nil
}

func (s *serverApi) DownloadFile(
	in *proto.FileRequest,
	str grpc.ServerStreamingServer[proto.FileResponse],
) error {

	return nil
}
