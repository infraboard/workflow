package pipeline

import (
	"fmt"

	"github.com/infraboard/workflow/conf"
	"github.com/infraboard/workflow/version"
)

func PipeLineObjectKey(namespace, id string) string {
	return fmt.Sprintf("%s/%s/%s", EtcdPipelinePrefix(), namespace, id)
}

func EtcdPipelinePrefix() string {
	return fmt.Sprintf("%s/%s/pipelines", conf.C().Etcd.Prefix, version.ServiceName)
}

func EtcdStepPrefix() string {
	return fmt.Sprintf("%s/%s/steps", conf.C().Etcd.Prefix, version.ServiceName)
}

func EtcdActionPrefix() string {
	return fmt.Sprintf("%s/%s/actions", conf.C().Etcd.Prefix, version.ServiceName)
}
