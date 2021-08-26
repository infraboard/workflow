package pkg

import (
	"fmt"

	"github.com/infraboard/mcube/pb/http"
	"google.golang.org/grpc"

	"github.com/infraboard/workflow/api/pkg/action"
	"github.com/infraboard/workflow/api/pkg/application"
	"github.com/infraboard/workflow/api/pkg/deploy"
	"github.com/infraboard/workflow/api/pkg/pipeline"
	"github.com/infraboard/workflow/api/pkg/template"
)

var (
	// Example 服务
	Application application.ServiceServer
	Pipeline    pipeline.ServiceServer
	Action      action.ServiceServer
	Template    template.ServiceServer
	Deploy      deploy.ServiceServer
)

var (
	servers       []Service
	successLoaded []string

	entrySet = http.NewEntrySet()
)

// InitV1GRPCAPI 初始化GRPC服务
func InitV1GRPCAPI(server *grpc.Server) {
	application.RegisterServiceServer(server, Application)
	pipeline.RegisterServiceServer(server, Pipeline)
	action.RegisterServiceServer(server, Action)
	template.RegisterServiceServer(server, Template)
	deploy.RegisterServiceServer(server, Deploy)
}

// HTTPEntry todo
func HTTPEntry() *http.EntrySet {
	return entrySet
}

// GetPathEntry todo
func GetPathEntry(path string) *http.Entry {
	es := HTTPEntry()
	for i := range es.Items {
		if es.Items[i].Path == path {
			return es.Items[i]
		}
	}

	return nil
}

// LoadedService 查询加载成功的服务
func LoadedService() []string {
	return successLoaded
}
func addService(name string, svr Service) {
	servers = append(servers, svr)
	successLoaded = append(successLoaded, name)
}

// Service 注册上的服务必须实现的方法
type Service interface {
	Config() error
	HTTPEntry() *http.EntrySet
}

// RegistryService 服务实例注册
func RegistryService(name string, svr Service) {
	switch value := svr.(type) {
	case application.ServiceServer:
		if Application != nil {
			registryError(name)
		}
		Application = value
		addService(name, svr)
	case pipeline.ServiceServer:
		if Pipeline != nil {
			registryError(name)
		}
		Pipeline = value
		addService(name, svr)
	case action.ServiceServer:
		if Pipeline != nil {
			registryError(name)
		}
		Action = value
		addService(name, svr)
	case template.ServiceServer:
		if Template != nil {
			registryError(name)
		}
		Template = value
		addService(name, svr)
	case deploy.ServiceServer:
		if Deploy != nil {
			registryError(name)
		}
		Deploy = value
		addService(name, svr)
	default:
		panic(fmt.Sprintf("unknown service type %s", name))
	}
}

func registryError(name string) {
	panic("service " + name + " has registried")
}

// InitService 初始化所有服务
func InitService() error {
	for _, s := range servers {
		if err := s.Config(); err != nil {
			return err
		}
		entrySet.Merge(s.HTTPEntry())
	}
	return nil
}
