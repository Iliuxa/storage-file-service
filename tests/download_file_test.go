package tests

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	proto "storage-file-service/proto/gen"
	"storage-file-service/tests/suite"
	"testing"
)

func TestDownloadFile_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	fileId := createFile(ctx, st, t)

	stream, err := st.StorageClient.DownloadFile(ctx, &proto.FileRequest{
		FileId: fileId,
	})
	require.NoError(t, err)

	res, err := stream.Recv()
	require.NoError(t, err)
	fileInfo := res.GetInfo()
	assert.NotNil(t, fileInfo)

	imageData := bytes.Buffer{}

	for {
		res, err = stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			require.Fail(t, "cannot write chunk data: %v", err)
		}
		chunk := res.GetChunkData()

		_, err = imageData.Write(chunk)
		if err != nil {
			require.Fail(t, "cannot write chunk data: %v", err)
		}
	}

	exePath, err := os.Getwd()
	if err != nil {
		require.Fail(t, "failed to get executable path: %w", err)
	}

	file, err := os.Create(exePath + "/tmp/" + fileInfo.FileName)
	if err != nil {
		require.Fail(t, "cannot create image file: %w", err)
	}
	defer file.Close()

	_, err = imageData.WriteTo(file)
	if err != nil {
		require.Fail(t, "cannot write image to file: %w", err)
	}
}

func createFile(ctx context.Context, st *suite.Suite, t *testing.T) uint64 {
	file, err := os.Open("pictures/img.png")
	if err != nil {
		require.Fail(t, "cannot open image file: ", err)
	}
	defer file.Close()

	stream, err := st.StorageClient.UploadFile(ctx)
	if err != nil {
		require.Fail(t, "cannot upload image: ", err)
	}

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		require.Fail(t, "cannot hash image")
	}

	req := &proto.UploadRequest{
		Data: &proto.UploadRequest_Info{
			Info: &proto.ImageInfo{
				FileName: "img.png",
				FileHash: fmt.Sprintf("%x", hasher.Sum(nil)),
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		require.Fail(t, "cannot send image info to server: ", err, stream.RecvMsg(nil))
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		require.Fail(t, "cannot seek file: ", err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			require.Fail(t, "cannot read chunk to buffer: ", err)
		}

		req = &proto.UploadRequest{
			Data: &proto.UploadRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			require.Fail(t, "cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		require.Fail(t, "cannot receive response: ", err)
	}

	return res.GetId()
}
