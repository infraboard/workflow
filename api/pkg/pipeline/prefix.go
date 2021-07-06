package pipeline

import (
	"fmt"

	"github.com/infraboard/workflow/version"
)

// func PipeLineObjectKey(namespace, id string) string {
// 	return fmt.Sprintf("%s/%s/%s", EtcdPipelinePrefix(), namespace, id)
// }
func PipeLineObjectKey(namespace, id string) string {
	return fmt.Sprintf("%s/%s", EtcdPipelinePrefix(), id)
}

func EtcdPipelinePrefix() string {
	return fmt.Sprintf("workflow/%s/%s/pipelines", version.OrgName, version.ServiceName)
}

func EtcdStepPrefix() string {
	return fmt.Sprintf("workflow/%s/%s/steps", version.OrgName, version.ServiceName)
}

func EtcdActionPrefix() string {
	return fmt.Sprintf("workflow/%s/%s/actions", version.OrgName, version.ServiceName)
}
