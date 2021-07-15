package file_test

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/workflow/node/controller/step/store/file"
)

func TestUpload(t *testing.T) {
	should := assert.New(t)

	buffer := ioutil.NopCloser(strings.NewReader("hello world"))
	store := file.NewStore()
	id, err := store.CreateObject(context.Background(),
		"c16mhsddrei91m4ri0jg.c3iqcama0brimaq08e40.2.1")
	if should.NoError(err) {
		t.Log("log object id: ", id)
		should.NoError(store.Upload(context.Background(), id, buffer))
	}

}
