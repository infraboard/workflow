package http

import (
	"fmt"
	"net/http"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/workflow/api/pkg/scm"
	"github.com/infraboard/workflow/api/pkg/scm/gitlab"
)

const (
	GitlabEventHeaderKey = "X-Gitlab-Event"
	GitlabEventTokenKey  = "X-Gitlab-Token"
)

func (h *handler) QuerySCMProject(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	srcType := qs.Get("scm_type")

	var (
		ps  *scm.ProjectSet
		err error
	)
	switch srcType {
	case "gitlab", "":
		repo := gitlab.NewSCM(qs.Get("scm_addr"), qs.Get("token"))
		ps, err = repo.ListProjects()
	case "github":
	case "coding":
	default:
		response.Failed(w, exception.NewBadRequest("unknown scm_type %s", srcType))
		return
	}

	if err != nil {
		response.Failed(w, err)
		return
	}

	response.Success(w, ps.Items)
}

func (h *handler) GitLabHookHanler(w http.ResponseWriter, r *http.Request) {
	eventType := r.Header.Get(GitlabEventHeaderKey)
	appID := r.Header.Get(GitlabEventTokenKey)
	fmt.Println(appID)
	switch eventType {
	case "Push Hook":
		event := scm.NewDefaultWebHookEvent()
		if err := request.GetDataFromRequest(r, event); err != nil {
			response.Failed(w, err)
			return
		}
		fmt.Println(event)
	default:
		return
	}
}
