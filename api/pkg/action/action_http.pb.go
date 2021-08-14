// Code generated by protoc-gen-go-http. DO NOT EDIT.

package action

import (
	http "github.com/infraboard/mcube/pb/http"
)

// HttpEntry todo
func HttpEntry() *http.EntrySet {
	set := &http.EntrySet{
		Items: []*http.Entry{
			{
				Path:         "/workflow.action.Service/CreateAction",
				FunctionName: "CreateAction",
			},
			{
				Path:         "/workflow.action.Service/QueryAction",
				FunctionName: "QueryAction",
			},
			{
				Path:         "/workflow.action.Service/DescribeAction",
				FunctionName: "DescribeAction",
			},
			{
				Path:         "/workflow.action.Service/UpdateAction",
				FunctionName: "UpdateAction",
			},
			{
				Path:         "/workflow.action.Service/DeleteAction",
				FunctionName: "DeleteAction",
			},
		},
	}
	return set
}
