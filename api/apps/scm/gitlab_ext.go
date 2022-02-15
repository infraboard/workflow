package scm

import (
	"fmt"
	"path"
)

func NewDefaultWebHookEvent() *WebHookEvent {
	return &WebHookEvent{
		Commits: []*Commit{},
	}
}

func NewProjectSet() *ProjectSet {
	return &ProjectSet{
		Items: []*Project{},
	}
}

func (e *WebHookEvent) ShortDesc() string {
	return fmt.Sprintf("%s %s [%s]", e.Ref, e.EventName, e.ObjectKind)
}

func (e *WebHookEvent) GetBranche() string {
	return path.Base(e.GetRef())
}
