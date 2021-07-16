package file

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

func NewUploader(id string) *Uploader {
	return &Uploader{
		id:     id,
		root:   "runner_log",
		parent: dateDir(),
	}
}

type Uploader struct {
	id     string
	root   string
	parent string
}

func (u *Uploader) DriverName() string {
	return "local_file"
}
func (u *Uploader) ObjectID() string {
	return path.Join(u.root, u.parent, u.id)
}
func (u *Uploader) Upload(ctx context.Context, stream io.ReadCloser) error {
	defer stream.Close()

	f, err := u.createFile(ctx)
	if err != nil {
		return fmt.Errorf("create file error, %s", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.ReadFrom(stream)
	if err != nil {
		return err
	}
	if err := w.Flush(); err != nil {
		return fmt.Errorf("flush file error, %s", err)
	}
	return nil
}

func (u *Uploader) createFile(ctx context.Context) (*os.File, error) {
	fp := u.ObjectID()
	if checkFileIsExist(fp) {
		return os.OpenFile(fp, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	}

	if err := os.MkdirAll(path.Dir(fp), os.ModePerm); err != nil {
		return nil, err
	}

	return os.Create(fp)
}

func dateDir() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%d/%d/%d", year, int(month), day)
}

// 判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(filepath string) bool {
	var exist = true
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
