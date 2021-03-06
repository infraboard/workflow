package impl

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/infraboard/workflow/api/apps/deploy"
)

func newQueryApplicationDeployRequest(req *deploy.QueryApplicationDeployRequest) *queryRequest {
	return &queryRequest{
		QueryApplicationDeployRequest: req,
	}
}

type queryRequest struct {
	*deploy.QueryApplicationDeployRequest
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
	filter["domain"] = r.Domain

	if r.Namespace != "" {
		filter["namespace"] = r.Namespace
	}
	if r.AppId != "" {
		filter["app_id"] = r.AppId
	}
	if r.Environment != "" {
		filter["environment"] = r.Environment
	}

	return filter
}

func newDescRequest(req *deploy.DescribeApplicationDeployRequest) (*describeRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &describeRequest{req}, nil
}

type describeRequest struct {
	*deploy.DescribeApplicationDeployRequest
}

func (req *describeRequest) FindFilter() bson.M {
	filter := bson.M{}

	if req.Namespace != "" {
		filter["namespace"] = req.Namespace
	}

	if req.Id != "" {
		filter["_id"] = req.Id
	}

	return filter
}
