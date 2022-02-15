package cmd

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/infraboard/mcube/types/ftime"
	"github.com/spf13/cobra"

	"github.com/infraboard/workflow/api/apps/node"
	"github.com/infraboard/workflow/api/client"
	"github.com/infraboard/workflow/conf"
	"github.com/infraboard/workflow/version"

	etcd_register "github.com/infraboard/workflow/api/apps/node/etcd"
	informer "github.com/infraboard/workflow/common/informers/step"
	si_impl "github.com/infraboard/workflow/common/informers/step/etcd"
	"github.com/infraboard/workflow/node/controller/step"
	controller "github.com/infraboard/workflow/node/controller/step"
)

var (
	// pusher service config option
	confType string
	confFile string
	confEtcd string
)

// startCmd represents the start command
var serviceCmd = &cobra.Command{
	Use:   "start",
	Short: "workflow-node 流水线node",
	Long:  `workflow-node 流水线node`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 初始化全局变量
		if err := loadGloabl(confType); err != nil {
			return err
		}
		cfg := conf.C()

		// 初始化全局client
		if err := loadGRPCClient(cfg); err != nil {
			return err
		}

		// 启动服务
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

		// 初始化服务
		svr, err := newService(cfg)
		if err != nil {
			return err
		}

		// 注册服务
		r, err := etcd_register.NewEtcdRegister(svr.node)
		if err != nil {
			svr.log.Warn(err)
		}
		r.Debug(zap.L().Named("Register"))
		defer r.UnRegiste()
		if err := r.Registe(); err != nil {
			return err
		}

		// 等待信号处理
		go svr.waitSign(ch)

		// 启动服务
		if err := svr.start(); err != nil {
			return err
		}
		return nil
	},
}

type service struct {
	info informer.Informer
	log  logger.Logger
	stop context.CancelFunc
	node *node.Node
	ctl  *step.Controller
}

func newService(cfg *conf.Config) (*service, error) {
	// Controller 实例
	rn := MakeRegistryNode(cfg)

	// 实例化Informer
	info := si_impl.NewFilterInformer(cfg.Etcd.GetClient(), informer.NewNodeFilter(rn))

	ctl := controller.NewController(rn.Name(), info, client.C())
	ctl.Debug(zap.L().Named("Node"))

	svr := &service{
		info: info,
		log:  zap.L().Named("CLI"),
		node: rn,
		ctl:  ctl,
	}
	return svr, nil
}

func (s *service) start() error {
	// 启动informer, Informer 需要先与Controller启动,避免事件丢失
	ctx, cancel := context.WithCancel(context.Background())
	s.stop = cancel
	defer cancel()

	if err := s.info.Watcher().Run(ctx); err != nil {
		s.log.Error(err)
	}

	// 启动controller
	if err := s.ctl.Run(ctx); err != nil {
		return err
	}

	return nil
}

// 占不做信号的具体区别
func (s *service) waitSign(sign chan os.Signal) {
	for {
		select {
		case sg := <-sign:
			switch v := sg.(type) {
			default:
				s.log.Infof("receive signal '%s', start graceful shutdown ...", v.String())
				s.stop()
				// 停止 总线
				s.log.Info("workflow node service stoped.")
				return
			}
		}
	}
}

func RandString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func MakeRegistryNode(cfg *conf.Config) *node.Node {
	hn, _ := os.Hostname()
	return &node.Node{
		InstanceName: hn,
		ServiceName:  version.ServiceName,
		Type:         node.NodeType,
		Address:      cfg.HTTP.Host,
		GitBranch:    version.GIT_BRANCH,
		GitCommit:    version.GIT_COMMIT,
		BuildEnv:     version.GO_VERSION,
		BuildAt:      version.BUILD_TIME,
		Online:       ftime.Now().Timestamp(),
		Prefix:       cfg.Etcd.Prefix,
		TTL:          cfg.Etcd.InstanceTTL,
		Interval:     time.Duration(cfg.Etcd.InstanceTTL/3) * time.Second,
	}
}

// config 为全局变量, 只需要load 即可全局可用户
// 日志需要初始化并配置
func loadGloabl(configType string) error {
	// 配置加载
	switch configType {
	case "file":
		err := conf.LoadConfigFromToml(confFile)
		if err != nil {
			return err
		}
	case "env":
		return errors.New("not implemented")
	case "etcd":
		return errors.New("not implemented")
	default:
		return errors.New("unknown config type")
	}
	// 加载日志组件
	lc := conf.C().Log
	var (
		logInitMsg string
		level      zap.Level
	)
	lv, err := zap.NewLevel(lc.Level)
	if err != nil {
		logInitMsg = fmt.Sprintf("%s, use default level INFO", err)
		level = zap.InfoLevel
	} else {
		level = lv
		logInitMsg = fmt.Sprintf("log level: %s", lv)
	}

	// 设置日志输出格式
	switch lc.Format {
	case conf.JSONFormat:
		err = zap.DevelopmentSetup(zap.WithLevel(level), zap.AsJSON())
	default:
		err = zap.DevelopmentSetup(zap.WithLevel(level))
	}
	if err != nil {
		return err
	}

	// 设置日志输出位置
	switch lc.To {
	case conf.ToFile:
		logconf := zap.DefaultConfig()
		logconf.Files.Name = "node.log"
		logconf.Files.Path = lc.PathDir
		logconf.Level = level
		if err := zap.Configure(logconf); err != nil {
			return err
		}
	}
	zap.L().Named("Init").Info(logInitMsg)
	return nil
}

// InitGRPCClient 初始化grpc客户端
func loadGRPCClient(cfg *conf.Config) error {
	cf := client.NewDefaultConfig()
	cf.SetAddress(cfg.GRPC.Addr())
	cf.SetClientCredentials(cfg.Keyauth.ClientID, cfg.Keyauth.ClientSecret)
	cli, err := client.NewClient(cf)
	if err != nil {
		return err
	}

	client.SetGlobal(cli)

	return err
}

func init() {
	serviceCmd.Flags().StringVarP(&confType, "config-type", "t", "file", "the service config type [file/env/etcd]")
	serviceCmd.Flags().StringVarP(&confFile, "config-file", "f", "etc/workflow.toml", "the service config from file")
	serviceCmd.Flags().StringVarP(&confEtcd, "config-etcd", "e", "127.0.0.1:2379", "the service config from etcd")
	RootCmd.AddCommand(serviceCmd)
}
