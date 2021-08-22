package gitlab_test

import (
	"fmt"
	"testing"

	"github.com/infraboard/workflow/common/repo/gitlab"
	"github.com/stretchr/testify/assert"
)

var (
	GitLabAddr    = "https://gitlab.com"
	PraviateToken = ""

	ProjectID int64 = 29032549
)

func TestListProject(t *testing.T) {
	should := assert.New(t)

	repo := gitlab.NewRepository(GitLabAddr, PraviateToken)
	ps, err := repo.ListProjects()
	should.NoError(err)
	fmt.Println(ps)
}

func TestAddProjectHook(t *testing.T) {
	should := assert.New(t)

	repo := gitlab.NewRepository(GitLabAddr, PraviateToken)

	req := gitlab.NewAddProjectHookRequest()
	req.ProjectID = ProjectID
	req.Hook = &gitlab.WebHook{
		PushEvents:          true,
		TagPushEvents:       true,
		MergeRequestsEvents: true,
		Token:               "9999",
		Url:                 "http://www.baidu.com",
	}
	ins, err := repo.AddProjectHook(req)
	should.NoError(err)
	fmt.Println(ins)
}

func TestDeleteProjectHook(t *testing.T) {
	should := assert.New(t)

	repo := gitlab.NewRepository(GitLabAddr, PraviateToken)

	req := gitlab.NewDeleteProjectReqeust(ProjectID, 8439846)
	err := repo.DeleteProjectHook(req)
	should.NoError(err)
}
