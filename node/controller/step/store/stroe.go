package store

import (
	"context"
	"io"
)

// 保存Runner运行中的日志
type WatcherOSS interface {
	StoreType() string
	Uploader
	Watcher
	Reader
}

// 用于上传日志
type Uploader interface {
	CreateObject(ctx context.Context, key string) (objectID string, err error)
	Upload(ctx context.Context, objectID string, steam io.ReadCloser) error
}

// 用于实时日志读取
type Watcher interface {
	Watch(ctx context.Context, objectID string, steam io.WriteCloser) error
}

// 用于历史日志访问
type Reader interface {
	ReadLine(ctx context.Context, objectID string, offset int64, total int64)
}
