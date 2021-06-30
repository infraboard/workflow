package version

import (
	"fmt"
)

var (
	// 组织的名称
	OrgName = "inforboard"
	// ServiceName 服务名称
	ServiceName = "workflow"
)

var (
	GIT_TAG    string
	GIT_COMMIT string
	GIT_BRANCH string
	BUILD_TIME string
	GO_VERSION string
)

// FullVersion show the version info
func FullVersion() string {
	version := fmt.Sprintf("Version   : %s\nBuild Time: %s\nGit Branch: %s\nGit Commit: %s\nGo Version: %s\n", GIT_TAG, BUILD_TIME, GIT_BRANCH, GIT_COMMIT, GO_VERSION)
	return version
}

// Short 版本缩写
func Short() string {
	short := ""
	if len(GIT_COMMIT) > 8 {
		short = GIT_COMMIT[:8]
	}
	return fmt.Sprintf("%s[%s %s]", GIT_TAG, BUILD_TIME, short)
}
