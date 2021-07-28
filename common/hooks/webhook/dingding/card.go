package dingding

type ActionCard struct {
	Title              string    `json:"title"`
	Text               string    `json:"text"`
	Buttons            []*Button `json:"btns"`
	ButtonsOrientation string    `json:"btnOrientation"`
	SingleTitle        string    `json:"singleTitle"`
	SingleURL          string    `json:"singleURL"`
}

type Button struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}
