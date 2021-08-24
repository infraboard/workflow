package gitlab

func NewDefaultWebHookEvent() *WebHookEvent {
	return &WebHookEvent{
		Commits: []*Commit{},
	}
}

type WebHookEvent struct {
	ObjectKind string    `json:"object_kind"`
	EventName  string    `json:"event_name"`
	Ref        string    `json:"ref"`
	UserID     int64     `json:"user_id"`
	Username   string    `json:"user_name"`
	UserAvatar string    `json:"user_avatar"`
	Project    *Project  `json:"project"`
	Commits    []*Commit `json:"commits"`
}

type Commit struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Title     string   `json:"title"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
	Author    *Author  `json:"author"`
	Added     []string `json:"added"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
