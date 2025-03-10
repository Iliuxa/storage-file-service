package storagegrpc

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"os"
	"storage-file-service/internal/domain"
	"storage-file-service/internal/usecase"
	proto "storage-file-service/proto/gen"
)

type serverApi struct {
	proto.UnimplementedFileServiceServer
	storage           usecase.StorageUsecase
	downloadSemaphore chan struct{}
	listSemaphore     chan struct{}
}

func Register(gRPCServer *grpc.Server, storage usecase.StorageUsecase, listLimit, downloadUploadLimit int) {
	proto.RegisterFileServiceServer(
		gRPCServer,
		&serverApi{
			storage:           storage,
			downloadSemaphore: make(chan struct{}, listLimit),
			listSemaphore:     make(chan struct{}, downloadUploadLimit),
		},
	)
}

func (s *serverApi) GetFileList(
	ctx context.Context,
	in *proto.FileListRequest,
) (*proto.FileListResponse, error) {
	s.listSemaphore <- struct{}{}
	defer func() { <-s.listSemaphore }()

	files, err := s.storage.GetFileList(ctx)
	if err != nil {
		if errors.Is(err, domain.ErrFileNotFound) {
			return nil, status.Error(codes.Internal, domain.ErrFileNotFound.Error())
		}
		return nil, status.Error(codes.Internal, "Failed to getting file")
	}
	return &proto.FileListResponse{Files: files}, nil
}

func (s *serverApi) UploadFile(
	stream grpc.ClientStreamingServer[proto.UploadRequest, proto.UploadResponse],
) error {
	s.downloadSemaphore <- struct{}{}
	defer func() { <-s.downloadSemaphore }()

	req, err := stream.Recv()
	if err != nil {
		return status.Error(codes.Internal, "failed to getting package")
	}
	fileName := req.GetInfo().GetFileName()
	fileHash := req.GetInfo().GetFileHash()

	imageData := bytes.Buffer{}

	for {
		err = contextError(stream.Context())
		if err != nil {
			return status.Error(codes.Internal, "failed to getting package")
		}

		req, err = stream.Recv()
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
		return status.Error(codes.Internal, "failed to getting file")
	}

	err = stream.SendAndClose(&proto.UploadResponse{Id: uint64(fileId)})
	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}
	return nil
}

func (s *serverApi) DownloadFile(
	in *proto.FileRequest,
	stream grpc.ServerStreamingServer[proto.FileResponse],
) error {
	s.downloadSemaphore <- struct{}{}
	defer func() { <-s.downloadSemaphore }()

	fileInfo, filePath, err := s.storage.DownloadFile(stream.Context(), uint(in.GetFileId()))
	if err != nil {
		return status.Error(codes.Internal, "failed to download file")
	}

	err = stream.Send(&proto.FileResponse{Data: fileInfo})
	if err != nil {
		return status.Error(codes.Internal, "failed to download file")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return status.Error(codes.Internal, "failed to download file")
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		err = contextError(stream.Context())
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot read chunck to buffer: %v", err)
		}

		res := &proto.FileResponse{
			Data: &proto.FileResponse_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send chunk to the client: %v", err)
		}
	}

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
