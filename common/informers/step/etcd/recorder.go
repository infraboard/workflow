package etcd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/types/ftime"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/api/apps/pipeline"
)

type recorder struct {
	log    logger.Logger
	client clientv3.KV
}

func (l *recorder) Update(step *pipeline.Step) error {
	step.UpdateAt = ftime.Now().Timestamp()
	objKey := pipeline.StepObjectKey(step.Key)
	objValue, err := json.Marshal(step)
	if err != nil {
		return err
	}

	l.log.Debugf("update step %s status %s %s ...", objKey, step.Status, string(objValue))
	if _, err := l.client.Put(context.Background(), objKey, string(objValue)); err != nil {
		return fmt.Errorf("update pipeline step '%s' to etcd3 failed: %s", objKey, err.Error())
	}
	return nil
}
