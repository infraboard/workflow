// Code generated by protoc-gen-go-http. DO NOT EDIT.

package application

import (
	http "github.com/infraboard/mcube/pb/http"
)

// HttpEntry todo
func HttpEntry() *http.EntrySet {
	set := &http.EntrySet{
		Items: []*http.Entry{
			{
				Path:         "/workflow.application.Service/CreateApplication",
				FunctionName: "CreateApplication",
			},
			{
				Path:         "/workflow.application.Service/UpdateApplication",
				FunctionName: "UpdateApplication",
			},
			{
				Path:         "/workflow.application.Service/QueryApplication",
				FunctionName: "QueryApplication",
			},
			{
				Path:         "/workflow.application.Service/DescribeApplication",
				FunctionName: "DescribeApplication",
			},
			{
				Path:         "/workflow.application.Service/DeleteApplication",
				FunctionName: "DeleteApplication",
			},
			{
				Path:         "/workflow.application.Service/HandleApplicationEvent",
				FunctionName: "HandleApplicationEvent",
			},
		},
	}
	return set
}
