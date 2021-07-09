package step

import (
	"fmt"
	"strings"

	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/common/cache"
)

// DefaultStoreIndexFunc todo
func DefaultStoreIndexFunc(obj interface{}) ([]string, error) {
	return []string{obj.(*pipeline.Step).MakeObjectKey()}, nil
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
	return obj.(*pipeline.Step).MakeObjectKey(), nil
}

// SplitMetaNamespaceKey returns the namespace and name that
// MetaNamespaceKeyFunc encoded into key.
//
// TODO: replace key-as-string with a key-as-struct so that this
// packing/unpacking won't be necessary.
func SplitMetaNamespaceKey(key string) (region, provider, excutor, id string, err error) {
	parts := strings.Split(key, "/")
	switch len(parts) {
	case 4:
		return parts[2], parts[3], "", "", nil
	case 5:
		return parts[2], parts[3], parts[4], "", nil
	case 6:
		return parts[2], parts[3], parts[4], parts[5], nil
	}
	return "", "", "", "", fmt.Errorf("unexpected key format: %q", key)
}
