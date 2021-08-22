package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
)

type Project struct {
	ID            int64  `json:"id"`
	Desc          string `json:"description"`
	Name          string `json:"name"`
	SshURL        string `json:"ssh_url_to_repo"`
	HttpURL       string `json:"http_url_to_repo"`
	NamespacePath string `json:"path_with_namespace"`
	HasSynced     bool   `json:"has_synced"`
}

func NewProjectSet() *ProjectSet {
	return &ProjectSet{
		Items: []*Project{},
	}
}

type ProjectSet struct {
	Items []*Project `json:"items"`
}

// https://gitlab.com/api/v4/projects?owned=true
// https://docs.gitlab.com/ce/api/projects.html
func (r *Repository) ListProjects() (*ProjectSet, error) {
	projectURL := r.resourceURL("projects", map[string]string{"owned": "true", "simple": "true"})
	req, err := r.newJSONRequest("GET", projectURL)
	if err != nil {
		return nil, err
	}

	// 发起请求
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取body
	bytesB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respString := string(bytesB)

	if (resp.StatusCode / 100) != 2 {
		return nil, fmt.Errorf("status code[%d] is not 200, response %s", resp.StatusCode, respString)
	}

	set := NewProjectSet()
	if err := json.Unmarshal(bytesB, &set.Items); err != nil {
		return nil, err
	}

	return set, nil
}

type WebHook struct {
	PushEventsBranchFilter string `json:"push_events_branch_filter"`
	PushEvents             bool   `json:"push_events"`
	TagPushEvents          bool   `json:"tag_push_events"`
	MergeRequestsEvents    bool   `json:"merge_requests_events"`
	Token                  string `json:"token"`
	Url                    string `json:"url"`
}

func (req *WebHook) FormValue() url.Values {
	val := make(url.Values)
	val.Set("push_events", fmt.Sprintf("%t", req.PushEvents))
	val.Set("tag_push_events", fmt.Sprintf("%t", req.TagPushEvents))
	val.Set("merge_requests_events", fmt.Sprintf("%t", req.MergeRequestsEvents))
	val.Set("token", req.Token)
	val.Set("url", req.Url)
	return val
}

func NewAddProjectHookRequest() *AddProjectHookRequest {
	return &AddProjectHookRequest{
		Hook: &WebHook{},
	}
}

type AddProjectHookRequest struct {
	ProjectID int64
	Hook      *WebHook
}

func NewAddProjectHookResponse() *AddProjectHookResponse {
	return &AddProjectHookResponse{
		WebHook: &WebHook{},
	}
}

type AddProjectHookResponse struct {
	ID int64 `json:"id"`
	*WebHook
}

// POST /projects/:id/hooks
// https://docs.gitlab.com/ce/api/projects.html#add-project-hook
func (r *Repository) AddProjectHook(in *AddProjectHookRequest) (*AddProjectHookResponse, error) {
	addHookURL := r.resourceURL(fmt.Sprintf("projects/%d/hooks", in.ProjectID), nil)
	req, err := r.newFormReqeust("POST", addHookURL, strings.NewReader(in.Hook.FormValue().Encode()))
	if err != nil {
		return nil, err
	}

	// 发起请求
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取body
	bytesB, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	respString := string(bytesB)

	if (resp.StatusCode / 100) != 2 {
		return nil, fmt.Errorf("status code[%d] is not 200, response %s", resp.StatusCode, respString)
	}

	ins := NewAddProjectHookResponse()
	if err := json.Unmarshal(bytesB, &ins); err != nil {
		return nil, err
	}

	return ins, nil
}

func NewDeleteProjectReqeust(projectID, hookID int64) *DeleteProjectReqeust {
	return &DeleteProjectReqeust{
		ProjectID: projectID,
		HookID:    hookID,
	}
}

type DeleteProjectReqeust struct {
	ProjectID int64
	HookID    int64
}

// DELETE /projects/:id/hooks/:hook_id
// https://docs.gitlab.com/ce/api/projects.html#delete-project-hook
func (r *Repository) DeleteProjectHook(in *DeleteProjectReqeust) error {
	addHookURL := r.resourceURL(fmt.Sprintf("projects/%d/hooks/%d", in.ProjectID, in.HookID), nil)
	req, err := r.newFormReqeust("DELETE", addHookURL, nil)
	if err != nil {
		return err
	}

	// 发起请求
	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取body
	bytesB, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respString := string(bytesB)

	if (resp.StatusCode / 100) != 2 {
		return fmt.Errorf("status code[%d] is not 200, response %s", resp.StatusCode, respString)
	}

	return nil
}
