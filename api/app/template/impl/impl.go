package impl

import (
	"context"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"google.golang.org/grpc"

	"github.com/infraboard/workflow/api/app/template"
	"github.com/infraboard/workflow/conf"
)

var (
	// Service 服务实例
	svr = &impl{}
)

type impl struct {
	col *mongo.Collection
	log logger.Logger

	template.UnimplementedServiceServer
}

func (s *impl) Config() error {
	db := conf.C().Mongo.GetDB()
	dc := db.Collection("actions")

	indexs := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "name", Value: bsonx.Int32(-1)},
				{Key: "version", Value: bsonx.Int32(-1)},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bsonx.Doc{{Key: "create_at", Value: bsonx.Int32(-1)}},
		},
	}

	_, err := dc.Indexes().CreateMany(context.Background(), indexs)
	if err != nil {
		return err
	}

	s.col = dc
	s.log = zap.L().Named("Action")

	return nil
}

func (s *impl) Name() string {
	return template.AppName
}

func (s *impl) Registry(server *grpc.Server) {
	template.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
