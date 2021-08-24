package scm

import "fmt"

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
