package all

import (
	// 加载服务模块
	_ "github.com/infraboard/workflow/api/pkg/action/grpc"
	_ "github.com/infraboard/workflow/api/pkg/action/http"
	_ "github.com/infraboard/workflow/api/pkg/application/grpc"
	_ "github.com/infraboard/workflow/api/pkg/application/http"
	_ "github.com/infraboard/workflow/api/pkg/pipeline/grpc"
	_ "github.com/infraboard/workflow/api/pkg/pipeline/http"
	_ "github.com/infraboard/workflow/api/pkg/template/grpc"
	_ "github.com/infraboard/workflow/api/pkg/template/http"
	_ "github.com/infraboard/workflow/api/pkg/trigger/http"
)
