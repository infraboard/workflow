package node

import (
	"errors"
	"fmt"
	"strings"

	"github.com/infraboard/workflow/conf"
	"github.com/infraboard/workflow/version"
)

// ParseInstanceKey 解析key中的服务信息
func ParseInstanceKey(key string) (serviceName, instanceName string, serviceType Type, err error) {
	kl := strings.Split(key, "/")
	if len(kl) != 5 {
		err = errors.New("key format error, must like inforboard/workflow/service/node/node-dev")
		return
	}
	serviceName, serviceType, instanceName = kl[1], Type(kl[3]), kl[4]
	return
}

func NodeObjectKey(key string) string {
	return fmt.Sprintf("%s/%s", EtcdNodePrefix(), key)
}

func EtcdNodePrefix() string {
	return fmt.Sprintf("%s/%s/service", conf.C().Etcd.Prefix, version.ServiceName)
}

func EtcdNodePrefixWithType(t Type) string {
	return fmt.Sprintf("%s/%s/service/%s", conf.C().Etcd.Prefix, version.ServiceName, t)
}
