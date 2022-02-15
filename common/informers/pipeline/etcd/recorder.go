package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/logger"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/apps/pipeline"
)

type recorder struct {
	log    logger.Logger
	client clientv3.KV
}

func (l *recorder) Update(t *pipeline.Pipeline) error {
	if t == nil {
		return fmt.Errorf("update nil pipeline")
	}

	if l.client == nil {
		return fmt.Errorf("etcd client is nil")
	}

	objKey := t.EtcdObjectKey()
	objValue, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if _, err := l.client.Put(context.Background(), objKey, string(objValue)); err != nil {
		return fmt.Errorf("update pipeline task '%s' to etcd3 failed: %s", objKey, err.Error())
	}
	return nil
}
