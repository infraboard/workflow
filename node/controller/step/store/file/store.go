package file

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
)

func NewStore() *Store {
	return &Store{
		dataDir: "runner_log",
		log:     zap.L().Named("File Store"),
	}
}

type Store struct {
	dataDir string
	log     logger.Logger
}

func (s *Store) StoreType() string {
	return "local_file"
}

func (s *Store) SetDataDir(path string) {
	s.dataDir = path
}

func (s *Store) SetLogger(log logger.Logger) {
	s.log = log
}

func (s *Store) CreateObject(ctx context.Context, key string) (objectID string, err error) {
	id := path.Join(dateDir(), key)
	fp := s.filePath(id)

	if checkFileIsExist(fp) {
		return id, nil
	}

	if err := os.MkdirAll(path.Dir(fp), os.ModePerm); err != nil {
		return "", err
	}
	if _, err := os.Create(fp); err != nil {
		return "", err
	}
	return id, nil
}

func (s *Store) Upload(ctx context.Context, id string, stream io.ReadCloser) error {
	defer stream.Close()

	filepath := s.filePath(id)

	if !checkFileIsExist(filepath) {
		return fmt.Errorf("object id not exist")
	}

	// 同一个key 采用覆盖写入
	f, err := os.OpenFile(filepath, os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("open file error, %s", err)
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

	s.log.Debugf("upload file success, %s", filepath)
	return nil
}

func (s *Store) Watch(ctx context.Context, key string, steam io.WriteCloser) error {
	return nil
}

func (s *Store) ReadLine(ctx context.Context, objectID string, offset int64, total int64) {

}

func (s *Store) filePath(id string) string {
	return path.Join(s.dataDir, id)
}

// 判断文件是否存在  存在返回 true 不存在返回false
func checkFileIsExist(filepath string) bool {
	var exist = true
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func dateDir() string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%d/%d/%d", year, int(month), day)
}
