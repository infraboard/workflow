package pipeline

import (
	"fmt"

	"github.com/infraboard/workflow/version"
)

func EtcdPipelinePrefix(prefix string) string {
	return fmt.Sprintf("%s/%s/pipeline", prefix, version.ServiceName)
}

func EtcdStepPrefix(prefix string) string {
	return fmt.Sprintf("%s/%s/step", prefix, version.ServiceName)
}
