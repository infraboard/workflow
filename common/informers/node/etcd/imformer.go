package etcd

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/common/cache"
	"github.com/infraboard/workflow/common/informers/node"
)

// NewNodeInformer todo
func NewInformer(client *clientv3.Client, filter node.NodeFilterHandler) node.Informer {
	return &Informer{
		log:     zap.L().Named("Node"),
		client:  client,
		filter:  filter,
		indexer: cache.NewIndexer(node.MetaNamespaceKeyFunc, node.DefaultStoreIndexers()),
	}
}

// Informer todo
type Informer struct {
	log     logger.Logger
	client  *clientv3.Client
	shared  *shared
	lister  *lister
	indexer cache.Indexer
	filter  node.NodeFilterHandler
}

func (i *Informer) GetStore() cache.Store {
	return i.indexer
}

func (i *Informer) Debug(l logger.Logger) {
	i.log = l
	i.shared.log = l
	i.lister.log = l
}

func (i *Informer) Watcher() node.Watcher {
	if i.shared != nil {
		return i.shared
	}
	i.shared = &shared{
		log:     i.log.Named("Watcher"),
		client:  clientv3.NewWatcher(i.client),
		indexer: i.indexer,
		filter:  i.filter,
	}
	return i.shared
}

func (i *Informer) Lister() node.Lister {
	if i.lister != nil {
		return i.lister
	}
	i.lister = &lister{
		log:    i.log.Named("Lister"),
		client: clientv3.NewKV(i.client),
	}
	return i.lister
}
