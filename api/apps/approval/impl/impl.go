package impl

import (
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"

	"github.com/infraboard/workflow/api/apps/approval"
	"github.com/infraboard/workflow/conf"
)

var (
	// Service 服务实例
	svr = &service{}
)

type service struct {
	col *mongo.Collection
	log logger.Logger

	approval.UnimplementedServiceServer
}

func (s *service) Config() error {
	db := conf.C().Mongo.GetDB()
	dc := db.Collection("approval")

	s.col = dc
	s.log = zap.L().Named(s.Name())
	return nil
}

func (s *service) Name() string {
	return approval.AppName
}

func (s *service) Registry(server *grpc.Server) {
	approval.RegisterServiceServer(server, svr)
}

func init() {
	app.RegistryGrpcApp(svr)
}
