package node

import (
	"fmt"

	"github.com/infraboard/workflow/version"
)

func EtcdNodePrefix() string {
	return fmt.Sprintf("%s/%s/service", version.OrgName, version.ServiceName)
}

func EtcdNodePrefixWithType(t Type) string {
	return fmt.Sprintf("%s/%s/service/%s", version.OrgName, version.ServiceName, t)
}
