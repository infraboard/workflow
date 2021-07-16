package file_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/workflow/node/controller/step/store/file"
)

func TestUpload(t *testing.T) {
	should := assert.New(t)

	buffer := io.NopCloser(strings.NewReader("hello world"))
	store := file.NewUploader("c16mhsddrei91m4ri0jg.c3iqcama0brimaq08e40.2.1")
	err := store.Upload(context.Background(), buffer)
	should.NoError(err)
}
