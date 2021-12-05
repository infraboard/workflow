package pipeline

import (
	"github.com/infraboard/workflow/api/app/pipeline"
	"github.com/infraboard/workflow/common/cache"
)

// DefaultStoreIndexFunc todo
func DefaultStoreIndexFunc(obj interface{}) ([]string, error) {
	return []string{obj.(*pipeline.Pipeline).MakeObjectKey()}, nil
}

// DefaultStoreIndexers todo
func DefaultStoreIndexers() cache.Indexers {
	indexers := cache.Indexers{}
	indexers["by_val"] = DefaultStoreIndexFunc
	return indexers
}

// ExplicitKey can be passed to MetaNamespaceKeyFunc if you have the key for
// the object but not the object itself.
type ExplicitKey string

// MetaNamespaceKeyFunc is a convenient default KeyFunc which knows how to make
// keys for API objects which implement meta.Interface.
// The key uses the format <namespace>/<name> unless <namespace> is empty, then
// it's just <name>.
//
// TODO: replace key-as-string with a key-as-struct so that this
// packing/unpacking won't be necessary.
func MetaNamespaceKeyFunc(obj interface{}) (string, error) {
	if key, ok := obj.(ExplicitKey); ok {
		return string(key), nil
	}

	pl, ok := obj.(*pipeline.Pipeline)
	if ok {
		return pl.MakeObjectKey(), nil
	}

	return "", nil
}
