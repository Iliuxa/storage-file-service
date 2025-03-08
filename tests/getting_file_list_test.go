package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	proto "storage-file-service/proto/gen"
	"storage-file-service/tests/suite"
	"testing"
)

func TestGetFileList_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	respReg, err := st.StorageClient.GetFileList(ctx, &proto.FileListRequest{})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetFiles())
}
