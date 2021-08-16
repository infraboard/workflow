package grpc

import (
	"context"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/pb/http"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/x/bsonx"

	"github.com/infraboard/workflow/api/pkg"
	"github.com/infraboard/workflow/api/pkg/template"
	"github.com/infraboard/workflow/conf"
)

var (
	// Service 服务实例
	Service = &impl{}
)

type impl struct {
	col *mongo.Collection
	template.UnimplementedServiceServer

	log logger.Logger
}

func (s *impl) Config() error {
	db := conf.C().Mongo.GetDB()
	dc := db.Collection("template")

	indexs := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{Key: "create_at", Value: bsonx.Int32(-1)}},
		},
	}

	_, err := dc.Indexes().CreateMany(context.Background(), indexs)
	if err != nil {
		return err
	}

	s.col = dc
	s.log = zap.L().Named("Template")

	return nil
}

// HttpEntry todo
func (s *impl) HTTPEntry() *http.EntrySet {
	return template.HttpEntry()
}

func init() {
	pkg.RegistryService("template", Service)
}
