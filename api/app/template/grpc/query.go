package grpc

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/mcube/pb/resource"
	"github.com/infraboard/workflow/api/app/template"
)

func newQueryActionRequest(req *template.QueryTemplateRequest) *queryRequest {
	return &queryRequest{
		QueryTemplateRequest: req,
	}
}

type queryRequest struct {
	*template.QueryTemplateRequest
}

func (r *queryRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	opt := &options.FindOptions{
		Sort:  bson.D{{Key: "create_at", Value: -1}},
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

func (r *queryRequest) FindFilter() bson.M {
	filter := bson.M{}

	cond1 := bson.M{}
	if r.Namespace != "" {
		cond1["namespace"] = r.Namespace
	}
	if r.Name != "" {
		cond1["name"] = r.Name
	}

	filter["$or"] = bson.A{
		cond1,
		bson.M{"visiable_mode": resource.VisiableMode_GLOBAL},
	}

	return filter
}

func newDescTemplateRequest(req *template.DescribeTemplateRequest) (*describeRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &describeRequest{req}, nil
}

type describeRequest struct {
	*template.DescribeTemplateRequest
}

func (req *describeRequest) FindFilter() bson.M {
	filter := bson.M{}

	filter["_id"] = req.Id

	return filter
}
