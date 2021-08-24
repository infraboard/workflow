package gitlab_test

import (
	"fmt"
	"testing"

	"github.com/infraboard/workflow/api/pkg/scm/gitlab"
	"github.com/stretchr/testify/assert"
)

var (
	GitLabAddr    = "https://gitlab.com"
	PraviateToken = ""

	ProjectID int64 = 29032549
)

func TestListProject(t *testing.T) {
	should := assert.New(t)

	repo := gitlab.NewSCM(GitLabAddr, PraviateToken)
	ps, err := repo.ListProjects()
	should.NoError(err)
	fmt.Println(ps)
}

func TestAddProjectHook(t *testing.T) {
	should := assert.New(t)

	repo := gitlab.NewSCM(GitLabAddr, PraviateToken)

	hook := &gitlab.WebHook{
		PushEvents:          true,
		TagPushEvents:       true,
		MergeRequestsEvents: true,
		Token:               "9999",
		Url:                 "http://www.baidu.com",
	}
	req := gitlab.NewAddProjectHookRequest(ProjectID, hook)
	ins, err := repo.AddProjectHook(req)
	should.NoError(err)
	fmt.Println(ins)
}

func TestDeleteProjectHook(t *testing.T) {
	should := assert.New(t)

	repo := gitlab.NewSCM(GitLabAddr, PraviateToken)

	req := gitlab.NewDeleteProjectReqeust(ProjectID, 8439846)
	err := repo.DeleteProjectHook(req)
	should.NoError(err)
}
