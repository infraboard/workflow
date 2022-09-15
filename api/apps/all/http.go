package all

import (
	// 加载服务模块
	_ "github.com/infraboard/workflow/api/apps/action/http"
	_ "github.com/infraboard/workflow/api/apps/pipeline/http"
	_ "github.com/infraboard/workflow/api/apps/template/http"
)
