package etcd

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	clientv3 "go.etcd.io/etcd/client/v3"

	"github.com/infraboard/workflow/common/cache"
	informer "github.com/infraboard/workflow/common/informers/step"
)

// NewSInformer todo
func NewSInformer(client *clientv3.Client) informer.Informer {
	return &Informer{
		log:     zap.L().Named("Step"),
		client:  client,
		indexer: cache.NewIndexer(informer.MetaNamespaceKeyFunc, informer.DefaultStoreIndexers()),
	}
}

// Informer todo
type Informer struct {
	log      logger.Logger
	client   *clientv3.Client
	shared   *shared
	lister   *lister
	recorder *recorder
	indexer  cache.Indexer
}

func (i *Informer) GetStore() cache.Store {
	return i.indexer
}

func (i *Informer) Debug(l logger.Logger) {
	i.log = l
	i.shared.log = l
	i.lister.log = l
}

func (i *Informer) Watcher() informer.Watcher {
	if i.shared != nil {
		return i.shared
	}
	i.shared = &shared{
		log:     i.log.Named("Watcher"),
		client:  clientv3.NewWatcher(i.client),
		indexer: i.indexer,
	}
	return i.shared
}

func (i *Informer) Lister() informer.Lister {
	if i.lister != nil {
		return i.lister
	}
	i.lister = &lister{
		log:    i.log.Named("Lister"),
		client: clientv3.NewKV(i.client),
	}
	return i.lister
}

func (i *Informer) Recorder() informer.Recorder {
	if i.recorder != nil {
		return i.recorder
	}
	i.recorder = &recorder{
		log:    i.log.Named("Recorder"),
		client: clientv3.NewKV(i.client),
	}
	return i.recorder
}
