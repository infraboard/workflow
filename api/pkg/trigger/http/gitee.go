package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/infraboard/mcube/http/response"
)

// GiteeTrigger 处理来自Gitee的WebHook事件
// Gitee WebHook 推送数据格式说明: https://gitee.com/help/articles/4186#article-header2
func (h *handler) GiteeTrigger(w http.ResponseWriter, r *http.Request) {
	event := NewGiteeWebHookEvent()
	event.ParseHeaderFromHTTP(r)

	h.log.Debugf("receive gitee event: %s", event.Event)

	// 读取body数据
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		response.Failed(w, fmt.Errorf("read body error, %s", err))
		return
	}

	// 是否需要挂载事件触发源信息

	// 解析事件(只取需要的部分)
	if err := json.Unmarshal(body, event.GiteeWebHookBody); err != nil {
		response.Failed(w, fmt.Errorf("unmarshal json body error, %s", err))
		return
	}

	// 根据URL, Clone镜像仓库
	response.Success(w, event)
}

func NewGiteeWebHookEvent() *GiteeWebHookEvent {
	return &GiteeWebHookEvent{
		&GiteeWebHookHeader{},
		&GiteeWebHookBody{},
	}
}

type GiteeWebHookEvent struct {
	*GiteeWebHookHeader
	*GiteeWebHookBody
}

func (e *GiteeWebHookEvent) HeadCommit() *Commit {
	if len(e.Commits) > 1 {
		return e.Commits[0]
	}

	return nil
}

type GiteeWebHookHeader struct {
	ContentType string
	UserAgent   string
	Token       string
	Timestamp   string
	Event       string
}

func (h *GiteeWebHookHeader) ParseHeaderFromHTTP(r *http.Request) {
	h.ContentType = r.Header.Get("Content-Type")
	h.UserAgent = r.UserAgent()
	h.Token = r.Header.Get("X-Gitee-Token")
	h.Timestamp = r.Header.Get("X-Gitee-Timestamp")
	h.Event = r.Header.Get("X-Gitee-Event")
}

type GiteeWebHookBody struct {
	Ref        string           `json:"ref"`
	Password   string           `json:"password"`
	Timestamp  int64            `json:"timestamp"`
	Sign       string           `json:"sign"`
	Repository *GiteeRepository `json:"repository"`
	Commits    []*Commit        `json:"commits"`
}

// GiteeRepository todo
type GiteeRepository struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	FullName   string `json:"full_name"`
	GitHttpUrl string `json:"git_http_url"`
	GitSshUrl  string `json:"git_ssh_url"`
}

type Commit struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Timestamp string    `json:"timestamp"`
	Url       string    `json:"url"`
	Added     []string  `json:"added"`
	Removed   []string  `json:"removed"`
	Modified  []string  `json:"modified"`
	Committer Committer `json:"committer"`
}

type Committer struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"user_name"`
	Url      string `json:"url"`
}
