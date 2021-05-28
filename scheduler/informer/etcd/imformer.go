package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/conf"
	"github.com/infraboard/workflow/scheduler/informer"
)

// NewScheduleInformer todo
func NewSchedulerInformer(cnf conf.Etcd) (informer.Informer, error) {
	if err := cnf.Validate(); err != nil {
		return nil, err
	}

	info := &PipelineInformer{cnf: cnf, log: zap.L().Named("Informer")}

	// 初始化客户端
	config := clientv3.Config{
		Endpoints:   cnf.Endpoints,
		DialTimeout: time.Duration(5) * time.Second,
		Username:    cnf.UserName,
		Password:    cnf.Password,
	}
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err = client.MemberList(ctx)
	if err != nil {
		return nil, fmt.Errorf("check etcd %s health by member list error, %s", cnf.Endpoints, err)
	}

	info.log.Debugf("connect to etcd %s success", cnf.Endpoints)
	info.client = client
	return info, nil
}

// PipelineInformer todo
type PipelineInformer struct {
	cnf    conf.Etcd
	log    logger.Logger
	client *clientv3.Client
	shared *shared
	lister *lister
}

func (i *PipelineInformer) Debug(l logger.Logger) {
	i.log = l
	i.shared.log = l
	i.lister.log = l
}

func (i *PipelineInformer) Watcher() informer.Watcher {
	if i.shared != nil {
		return i.shared
	}
	i.shared = &shared{
		log:    i.log.Named("Watcher"),
		client: clientv3.NewWatcher(i.client),
		prefix: i.cnf.Prefix,
	}
	return i.shared
}

func (i *PipelineInformer) Lister() informer.Lister {
	if i.lister != nil {
		return i.lister
	}
	i.lister = &lister{
		log:    i.log.Named("Lister"),
		client: clientv3.NewKV(i.client),
		prefix: i.cnf.Prefix,
	}
	return i.lister
}
