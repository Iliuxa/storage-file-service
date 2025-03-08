package tests

import (
	"bufio"
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

func TestUploadFile_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

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

	require.NoError(t, err)
	assert.NotEmpty(t, res.GetId())
}
