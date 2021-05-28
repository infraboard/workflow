package application

// NewApplication todo
func NewApplication(req *CreateApplicationRequest) *Application {
	return &Application{
		Id:   "mock id",
		Name: req.Name,
	}
}

// NewApplicationSet 实例
func NewApplicationSet() *ApplicationSet {
	return &ApplicationSet{
		Items: []*Application{},
	}
}
