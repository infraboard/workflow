package pipeline

import (
	"fmt"

	"github.com/infraboard/workflow/version"
)

func PipeLineObjectKey(namespace, id string) string {
	return fmt.Sprintf("%s/%s/%s", EtcdPipelinePrefix(), namespace, id)
}

func EtcdPipelinePrefix() string {
	return fmt.Sprintf("%s/%s/pipelines", version.OrgName, version.ServiceName)
}

func EtcdStepPrefix() string {
	return fmt.Sprintf("%s/%s/steps", version.OrgName, version.ServiceName)
}
