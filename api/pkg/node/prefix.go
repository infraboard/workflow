package node

import (
	"fmt"

	"github.com/infraboard/workflow/conf"
	"github.com/infraboard/workflow/version"
)

func EtcdNodePrefix() string {

	return fmt.Sprintf("%s/%s/service", conf.C().Etcd.Prefix, version.ServiceName)
}

func EtcdNodePrefixWithType(t Type) string {
	return fmt.Sprintf("%s/%s/service/%s", conf.C().Etcd.Prefix, version.ServiceName, t)
}
