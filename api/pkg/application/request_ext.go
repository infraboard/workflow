package application

import "github.com/infraboard/mcube/http/request"

// NewQueryBookRequest 查询book列表
func NewQueryBookRequest(page *request.PageRequest) *QueryApplicationRequest {
	return &QueryApplicationRequest{
		Page: &page.PageRequest,
	}
}
