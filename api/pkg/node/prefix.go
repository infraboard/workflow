package node

import (
	"fmt"

	"github.com/infraboard/workflow/version"
)

func EtcdNodePrefix(prefix string) string {
	return fmt.Sprintf("%s/%s/service", prefix, version.ServiceName)
}

func EtcdNodePrefixWithType(prefix string, t Type) string {
	return fmt.Sprintf("%s/%s/service/%s", prefix, version.ServiceName, t)
}
