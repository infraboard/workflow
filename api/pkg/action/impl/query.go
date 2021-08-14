package impl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/mcube/pb/resource"
	"github.com/infraboard/workflow/api/pkg/action"
)

func newQueryActionRequest(req *action.QueryActionRequest) *queryRequest {
	return &queryRequest{
		QueryActionRequest: req,
	}
}

type queryRequest struct {
	*action.QueryActionRequest
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
	return filter
}

func newDescActionRequest(req *action.DescribeActionRequest) (*describeRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &describeRequest{req}, nil
}

type describeRequest struct {
	*action.DescribeActionRequest
}

func (req *describeRequest) FindFilter() bson.M {
	filter := bson.M{}

	filter["name"] = req.Name
	filter["version"] = req.Version
	filter["$or"] = bson.A{
		bson.M{"namespace": req.Namespace},
		bson.M{"visiable_mode": resource.VisiableMode_GLOBAL.String()},
	}

	return filter
}

func newDeleteActionRequest(req *action.DeleteActionRequest) (*deleteRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &deleteRequest{req}, nil
}

type deleteRequest struct {
	*action.DeleteActionRequest
}

func (req *deleteRequest) DeleteFilter() bson.M {
	filter := bson.M{}

	filter["name"] = req.Name
	filter["version"] = req.Version
	filter["namespace"] = req.Namespace

	return filter
}
