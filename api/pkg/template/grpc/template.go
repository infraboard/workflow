package grpc

import (
	"context"

	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/infraboard/workflow/api/pkg/template"
)

func (i *impl) CreateTemplate(ctx context.Context, req *template.CreateTemplateRequest) (
	*template.Template, error) {
	a, err := template.NewTemplate(req)
	if err != nil {
		return nil, err
	}

	if _, err := i.col.InsertOne(context.TODO(), a); err != nil {
		return nil, exception.NewInternalServerError("inserted a template document error, %s", err)
	}
	return a, nil
}

func (i *impl) QueryTemplate(ctx context.Context, req *template.QueryTemplateRequest) (
	*template.TemplateSet, error) {
	query := newQueryActionRequest(req)
	resp, err := i.col.Find(context.TODO(), query.FindFilter(), query.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find template error, error is %s", err)
	}

	set := template.NewTemplateSet()
	// 循环
	for resp.Next(context.TODO()) {
		a := template.NewDefaultTemplate()
		if err := resp.Decode(a); err != nil {
			return nil, exception.NewInternalServerError("decode template error, error is %s", err)
		}

		set.Add(a)
	}

	// count
	count, err := i.col.CountDocuments(context.TODO(), query.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get template count error, error is %s", err)
	}
	set.Total = count
	return set, nil
}

func (i *impl) DescribeTemplate(ctx context.Context, req *template.DescribeTemplateRequest) (
	*template.Template, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate DescribeTemplateRequest error, %s", err)
	}

	desc, err := newDescTemplateRequest(req)
	if err != nil {
		return nil, exception.NewBadRequest(err.Error())
	}

	ins := template.NewDefaultTemplate()
	if err := i.col.FindOne(context.TODO(), desc.FindFilter()).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("template %s not found", req)
		}

		return nil, exception.NewInternalServerError("find template %s error, %s", req.Id, err)
	}

	return ins, nil
}

func (i *impl) UpdateAction(context.Context, *template.UpdateTemplateRequest) (
	*template.Template, error) {
	return nil, nil
}

func (i *impl) DeleteTemplate(ctx context.Context, req *template.DeleteTemplateRequest) (
	*template.Template, error) {
	ins, err := i.DescribeTemplate(ctx, template.NewDescribeTemplateRequestWithID(req.Id))
	if err != nil {
		return nil, err
	}

	if _, err := i.col.DeleteOne(context.TODO(), bson.M{"_id": req.Id}); err != nil {
		return nil, err
	}

	return ins, nil
}
