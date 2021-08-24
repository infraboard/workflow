package scm

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
