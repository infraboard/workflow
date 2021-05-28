package informer

// Informer 负责事件通知
type Informer interface {
	Watcher() Watcher
	Lister() Lister
}
