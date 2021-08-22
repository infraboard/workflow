package http

import (
	"net/http"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/workflow/common/repo/gitlab"
)

func (h *handler) QueryRepoProject(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	srcType := qs.Get("scm_type")

	var (
		ps  *gitlab.ProjectSet
		err error
	)
	switch srcType {
	case "gitlab", "":
		repo := gitlab.NewRepository(qs.Get("scm_addr"), qs.Get("token"))
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
