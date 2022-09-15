package client

import (
	kc "github.com/infraboard/keyauth/client"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"github.com/infraboard/workflow/api/apps/action"
	"github.com/infraboard/workflow/api/apps/pipeline"
	"github.com/infraboard/workflow/api/apps/template"
)

var (
	client *ClientSet
)

// SetGlobal todo
func SetGlobal(cli *ClientSet) {
	client = cli
}

// C Global
func C() *ClientSet {
	return client
}

// NewClient todo
func NewClientSet(conf *kc.Config) (*ClientSet, error) {
	zap.DevelopmentSetup()
	log := zap.L()

	conn, err := grpc.Dial(conf.Address(), grpc.WithInsecure(), grpc.WithPerRPCCredentials(conf.Authentication))
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		conn: conn,
		log:  log,
	}, nil
}

// Client 客户端
type ClientSet struct {
	conn *grpc.ClientConn
	log  logger.Logger
}

// Example todo
func (c *ClientSet) Pipeline() pipeline.ServiceClient {
	return pipeline.NewServiceClient(c.conn)
}

// Example todo
func (c *ClientSet) Action() action.ServiceClient {
	return action.NewServiceClient(c.conn)
}

// Example todo
func (c *ClientSet) Template() template.ServiceClient {
	return template.NewServiceClient(c.conn)
}
