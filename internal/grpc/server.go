package storagegrpc

import (
	"bytes"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
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
			return nil, status.Error(codes.Internal, "Failed to getting file")
		}
		return nil, status.Error(codes.Internal, "Failed to getting file")
	}
	return &proto.FileListResponse{Files: files}, nil
}

func (s *serverApi) UploadFile(
	stream grpc.ClientStreamingServer[proto.UploadRequest, proto.UploadResponse],
) error {
	req, err := stream.Recv()
	if err != nil {
		panic(333)
	}
	fileName := req.GetInfo().GetFileName()
	fileHash := req.GetInfo().GetFileHash()

	imageData := bytes.Buffer{}

	for {
		err := contextError(stream.Context())
		if err != nil {
			return err
		}

		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot receive chunk data: %v", err)
		}

		chunk := req.GetChunkData()

		_, err = imageData.Write(chunk)
		if err != nil {
			return status.Errorf(codes.Internal, "cannot write chunk data: %v", err)
		}
	}

	fileId, err := s.storage.UploadFile(stream.Context(), &imageData, fileName, fileHash)
	if err != nil {
		return status.Error(codes.Internal, "Failed to getting file")
	}

	err = stream.SendAndClose(&proto.UploadResponse{Id: uint64(fileId)})
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}
	return nil
}

func (s *serverApi) DownloadFile(
	in *proto.FileRequest,
	str grpc.ServerStreamingServer[proto.FileResponse],
) error {

	return nil
}

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return nil
	case context.DeadlineExceeded:
		return nil
	default:
		return nil
	}
}
