package node

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/infraboard/mcube/logger"
)

const (
	// API 提供API访问的服务
	APIType = Type("api")
	// Worker 后台作业服务
	NodeType = Type("node")
	// Scheduler 调度器
	SchedulerType = Type("scheduler")
)

type Type string

// ParseEtcdNode tdo
func LoadNodeFromBytes(value []byte) (*Node, error) {
	n := new(Node)
	n.Tag = map[string]string{}
	// 解析Value
	if len(value) > 0 {
		if err := json.Unmarshal(value, n); err != nil {
			return nil, fmt.Errorf("unmarshal node error, vaule(%s) %s", string(value), err)
		}
	}

	// 校验合法性
	if string(value) == "" {
		return nil, nil
	}

	if err := n.Validate(); err != nil {
		return nil, err
	}
	return n, nil
}

// Node todo
type Node struct {
	Region          string            `json:"region,omitempty"`
	ResourceVersion int64             `json:"resourceVersion,omitempty"`
	InstanceName    string            `json:"instance_name,omitempty"`
	ServiceName     string            `json:"service_name,omitempty"`
	Type            Type              `json:"type,omitempty"`
	Address         string            `json:"address,omitempty"`
	Version         string            `json:"version,omitempty"`
	GitBranch       string            `json:"git_branch,omitempty"`
	GitCommit       string            `json:"git_commit,omitempty"`
	BuildEnv        string            `json:"build_env,omitempty"`
	BuildAt         string            `json:"build_at,omitempty"`
	Online          int64             `json:"online,omitempty"`
	Tag             map[string]string `json:"tag,omitempty"`

	Prefix   string        `json:"-"`
	Interval time.Duration `json:"-"`
	TTL      int64         `json:"-"`
}

func (n *Node) Name() string {
	return fmt.Sprintf("%s.%s", n.ServiceName, n.InstanceName)
}

func (n *Node) Validate() error {
	if n.InstanceName == "" && n.ServiceName == "" || n.Type == "" {
		return errors.New("service instance name or service_name or type not config")
	}
	return nil
}

// MakeRegistryKey 构建etcd对应的key
// 例如: inforboard/workflow/service/node/node-01
func (n *Node) MakeRegistryKey() string {
	return fmt.Sprintf("%s/%s/service/%s/%s", n.Prefix, n.ServiceName, n.Type, n.InstanceName)
}

// ParseInstanceKey 解析key中的服务信息
func ParseInstanceKey(key string) (serviceName, instanceName string, serviceType Type, err error) {
	kl := strings.Split(key, "/")
	if len(kl) != 5 {
		err = errors.New("key format error, must like inforboard/workflow/service/node/node-dev")
		return
	}
	serviceName, serviceType, instanceName = kl[1], Type(kl[3]), kl[4]
	return
}

// Register 服务注册接口
type Register interface {
	Debug(logger.Logger)
	Registe() error
	UnRegiste() error
}

// HeatbeatResonse 心态的返回
type HeatbeatResonse interface {
	TTL() int64
}
