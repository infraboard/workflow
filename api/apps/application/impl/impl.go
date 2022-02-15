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

	"github.com/infraboard/workflow/api/apps/application"
	"github.com/infraboard/workflow/api/apps/pipeline"
	"github.com/infraboard/workflow/conf"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	col      *mongo.Collection
	log      logger.Logger
	pipeline pipeline.ServiceServer

	platform string

	application.UnimplementedServiceServer
}

func (s *service) Config() error {
	db := conf.C().Mongo.GetDB()
	dc := db.Collection("dev_application")

	indexs := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "namespace", Value: bsonx.Int32(-1)},
				{Key: "name", Value: bsonx.Int32(-1)},
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

	s.pipeline = app.GetGrpcApp(pipeline.AppName).(pipeline.ServiceServer)

	s.platform = conf.C().App.Platform
	s.col = dc
	s.log = zap.L().Named("Application")
	return nil
}

func (s *service) Name() string {
	return application.AppName
}

func (s *service) Registry(server *grpc.Server) {
	application.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
