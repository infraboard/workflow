package store

import (
	"context"
	"io"

	"github.com/infraboard/workflow/node/controller/step/store/file"
)

// 保存Runner运行中的日志
type StoreFactory interface {
	NewFileUpdater(key string) Uploader
}

// 用于上传日志
type Uploader interface {
	DriverName() string
	ObjectID() string
	Upload(ctx context.Context, steam io.ReadCloser) error
}

func NewStore() *Store {
	return &Store{}
}

type Store struct{}

func (s *Store) NewFileUpdater(key string) Uploader {
	return file.NewUploader(key)
}
