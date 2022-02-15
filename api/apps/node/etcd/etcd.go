package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/workflow/api/apps/node"
	"github.com/infraboard/workflow/conf"
)

type etcd struct {
	leaseID        clientv3.LeaseID
	client         *clientv3.Client
	requestTimeout time.Duration
	isStopped      bool
	instanceKey    string
	instanceValue  string
	stopInstance   chan struct{}
	keepAliveStop  context.CancelFunc
	node           *node.Node
	logger.Logger
}

// NewEtcdRegister 初始化一个基于etcd的实例注册器
func NewEtcdRegister(node *node.Node) (node.Register, error) {
	if err := node.Validate(); err != nil {
		return nil, err
	}

	etcdR := new(etcd)
	etcdR.client = conf.C().Etcd.GetClient()
	etcdR.stopInstance = make(chan struct{}, 1)
	etcdR.requestTimeout = time.Duration(5) * time.Second
	etcdR.node = node

	// 注册服务的key
	sjson, err := json.Marshal(node)
	if err != nil {
		return nil, err
	}
	etcdR.instanceValue = string(sjson)
	etcdR.instanceKey = node.MakeObjectKey()
	return etcdR, nil
}

// node use to registe serice endpoint to etcd. when etcd is down,
// node can retry to registe util the etcd up.
//
// name is service name, use to discovery service address, eg. keyauth
// host and port is service endpoint, eg. 127.0.0.0:50000
// target is etcd addr, eg. 127.0.0.0:2379
// interval is service refresh interval, eg. 10s
// ttl is service ttl, eg. 15
// TODO: 判断服务是否已经被其他人注册了, 如果注册了 则需要更换名称才能注册
func (e *etcd) Registe() error {
	// 后台续约
	// 并没有直接使用KeepAlive, 因为存在偶然端口, 就不续约的情况
	ctx, cancel := context.WithCancel(context.Background())
	e.keepAliveStop = cancel

	// 初始化注册
	if err := e.addOnce(); err != nil {
		e.Errorf("registry error, %s", err)
		return err
	}
	e.Infof("服务实例(%s)注册成功", e.instanceKey)

	// keep alive
	go e.keepAlive(ctx)
	return nil
}

func (e *etcd) Debug(log logger.Logger) {
	e.Logger = log
}

func (e *etcd) addOnce() error {
	// 获取leaseID
	resp, err := e.client.Lease.Grant(context.TODO(), e.node.TTL)
	if err != nil {
		return fmt.Errorf("get etcd lease id error, %s", err)
	}
	e.leaseID = resp.ID

	// 写入key
	if _, err := e.client.Put(context.Background(), e.instanceKey, e.instanceValue, clientv3.WithLease(e.leaseID)); err != nil {
		return fmt.Errorf("registe service '%s' with ttl to etcd3 failed: %s", e.instanceKey, err.Error())
	}
	e.instanceKey = e.instanceKey
	return nil
}

func (e *etcd) keepAlive(ctx context.Context) {
	// 不停的续约
	interval := e.node.TTL / 5
	e.Infof("keep alive lease interval is %d second", interval)
	tk := time.NewTicker(time.Duration(interval) * time.Second)
	defer tk.Stop()
	for {
		select {
		case <-ctx.Done():
			e.Infof("keepalive goroutine exit")
			return
		case <-tk.C:
			Opctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err := e.client.Lease.KeepAliveOnce(Opctx, e.leaseID)
			if err != nil {
				if strings.Contains(err.Error(), "requested lease not found") {
					// 避免程序卡顿造成leaseID失效(比如mac 电脑休眠))
					if err := e.addOnce(); err != nil {
						e.Errorf("refresh registry error, %s", err)
					} else {
						e.Warn("refresh registry success")
					}
				}
				e.Errorf("lease keep alive error, %s", err)
			} else {
				e.Debugf("lease keep alive key: %s", e.instanceKey)
			}
		}
	}
}

// UnRegiste delete nodeed service from etcd, if etcd is down
// unnode while timeout.
func (e *etcd) UnRegiste() error {
	if e.isStopped {
		return errors.New("the instance has unregisted")
	}
	// delete instance key
	e.stopInstance <- struct{}{}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if resp, err := e.client.Delete(ctx, e.instanceKey); err != nil {
		e.Warnf("unregiste '%s' failed: connect to etcd server timeout, %s", e.instanceKey, err.Error())
	} else {
		if resp.Deleted == 0 {
			e.Infof("unregiste '%s' failed, the key not exist", e.instanceKey)
		} else {
			e.Infof("服务实例(%s)注销成功", e.instanceKey)
		}
	}
	// revoke lease
	_, err := e.client.Lease.Revoke(context.TODO(), e.leaseID)
	if err != nil {
		e.Warnf("revoke lease error, %s", err)
		return err
	}
	e.isStopped = true
	// 停止续约心态
	e.keepAliveStop()
	return nil
}
