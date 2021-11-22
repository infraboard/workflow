package all

import (
	// 加载服务模块
	_ "github.com/infraboard/workflow/api/app/action/http"
	_ "github.com/infraboard/workflow/api/app/application/http"
	_ "github.com/infraboard/workflow/api/app/deploy/http"
	_ "github.com/infraboard/workflow/api/app/pipeline/http"
	_ "github.com/infraboard/workflow/api/app/template/http"
	_ "github.com/infraboard/workflow/api/app/trigger/http"
)
