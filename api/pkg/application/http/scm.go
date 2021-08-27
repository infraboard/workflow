package http

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/grpc/gcontext"
	"github.com/infraboard/mcube/http/request"
	"github.com/infraboard/mcube/http/response"

	"github.com/infraboard/workflow/api/pkg/application"
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
	switch eventType {
	case "Push Hook":
		event := scm.NewDefaultWebHookEvent()
		if err := request.GetDataFromRequest(r, event); err != nil {
			response.Failed(w, err)
			return
		}

		req := application.NewApplicationEvent(appID, event)
		h.log.Debugf("application %d accept event: %s", appID, event)

		var header, trailer metadata.MD
		_, err := h.service.HandleApplicationEvent(
			r.Context(),
			req,
			grpc.Header(&header),
			grpc.Trailer(&trailer),
		)
		if err != nil {
			response.Failed(w, gcontext.NewExceptionFromTrailer(trailer, err))
			return
		}
		response.Success(w, fmt.Sprintf("event %s has accept", event.ShortDesc()))
		return
	default:
		response.Failed(w, fmt.Errorf("known gitlab event type %s", eventType))
		return
	}
}
