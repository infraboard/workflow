package conf

import (
	"context"
	"fmt"
	"time"

	kc "github.com/infraboard/keyauth/client"

	"github.com/infraboard/mcube/bus/broker/nats"
	"github.com/infraboard/mcube/cache/memory"
	"github.com/infraboard/mcube/cache/redis"
	"github.com/infraboard/mcube/logger/zap"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mgoClient  *mongo.Client
	etcdClient *clientv3.Client
)

func newConfig() *Config {
	return &Config{
		App:     newDefaultAPP(),
		HTTP:    newDefaultHTTP(),
		GRPC:    newDefaultGRPC(),
		Log:     newDefaultLog(),
		Mongo:   newDefaultMongoDB(),
		Cache:   newDefaultCache(),
		Keyauth: newDefaultKeyauth(),
		Etcd:    newDefaultEtcd(),
		Nats:    nats.NewDefaultConfig(),
		Bus:     new(bus),
	}
}

// Config 应用配置
type Config struct {
	App     *app         `toml:"app"`
	HTTP    *http        `toml:"http"`
	GRPC    *grpc        `toml:"grpc"`
	Log     *log         `toml:"log"`
	Mongo   *mongodb     `toml:"mongodb"`
	Keyauth *keyauth     `toml:"keyauth"`
	Cache   *_cache      `toml:"cache"`
	Etcd    *Etcd        `toml:"etcd"`
	Nats    *nats.Config `toml:"nats"`
	Bus     *bus         `toml:"bus"`
}

type bus struct {
	Type string `toml:"type" env:"BUS_TYPE"`
}

type app struct {
	Name     string `toml:"name" env:"APP_NAME"`
	Key      string `toml:"key" env:"APP_KEY"`
	Platform string `toml:"platform" env:"APP_PLATFORM"`
}

func newDefaultAPP() *app {
	return &app{
		Name: "workflow",
		Key:  "default",
	}
}

type http struct {
	Host      string `toml:"host" env:"HTTP_HOST"`
	Port      string `toml:"port" env:"HTTP_PORT"`
	EnableSSL bool   `toml:"enable_ssl" env:"HTTP_ENABLE_SSL"`
	CertFile  string `toml:"cert_file" env:"HTTP_CERT_FILE"`
	KeyFile   string `toml:"key_file" env:"HTTP_KEY_FILE"`
}

func (a *http) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultHTTP() *http {
	return &http{
		Host: "127.0.0.1",
		Port: "8050",
	}
}

type grpc struct {
	Host      string `toml:"host" env:"GRPC_HOST"`
	Port      string `toml:"port" env:"GRPC_PORT"`
	EnableSSL bool   `toml:"enable_ssl" env:"GRPC_ENABLE_SSL"`
	CertFile  string `toml:"cert_file" env:"GRPC_CERT_FILE"`
	KeyFile   string `toml:"key_file" env:"GRPC_KEY_FILE"`
}

func (a *grpc) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultGRPC() *grpc {
	return &grpc{
		Host: "127.0.0.1",
		Port: "18050",
	}
}

type log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
}

func newDefaultLog() *log {
	return &log{
		Level:   "debug",
		PathDir: "logs",
		Format:  "text",
		To:      "stdout",
	}
}

// Auth auth 配置
type keyauth struct {
	Host         string `toml:"host" env:"KEYAUTH_HOST"`
	Port         string `toml:"port" env:"KEYAUTH_PORT"`
	ClientID     string `toml:"client_id" env:"KEYAUTH_CLIENT_ID"`
	ClientSecret string `toml:"client_secret" env:"KEYAUTH_CLIENT_SECRET"`
}

func (a *keyauth) Addr() string {
	return a.Host + ":" + a.Port
}

func (a *keyauth) Client() (*kc.Client, error) {
	if kc.C() == nil {
		conf := kc.NewDefaultConfig()
		conf.SetAddress(a.Addr())
		conf.SetClientCredentials(a.ClientID, a.ClientSecret)
		client, err := kc.NewClient(conf)
		if err != nil {
			return nil, err
		}
		kc.SetGlobal(client)
	}

	return kc.C(), nil
}

func newDefaultKeyauth() *keyauth {
	return &keyauth{}
}

func newDefaultMongoDB() *mongodb {
	return &mongodb{
		Database:  "",
		Endpoints: []string{"127.0.0.1:27017"},
	}
}

type mongodb struct {
	Endpoints []string `toml:"endpoints" env:"MONGO_ENDPOINTS" envSeparator:","`
	UserName  string   `toml:"username" env:"MONGO_USERNAME"`
	Password  string   `toml:"password" env:"MONGO_PASSWORD"`
	Database  string   `toml:"database" env:"MONGO_DATABASE"`
}

// Client 获取一个全局的mongodb客户端连接
func (m *mongodb) Client() *mongo.Client {
	if mgoClient == nil {
		panic("please load mongo client first")
	}

	return mgoClient
}

func (m *mongodb) GetDB() *mongo.Database {
	return m.Client().Database(m.Database)
}

func (m *mongodb) getClient() (*mongo.Client, error) {
	opts := options.Client()

	cred := options.Credential{
		AuthSource: m.Database,
	}

	if m.UserName != "" && m.Password != "" {
		cred.Username = m.UserName
		cred.Password = m.Password
		cred.PasswordSet = true
		opts.SetAuth(cred)
	}
	opts.SetHosts(m.Endpoints)
	opts.SetConnectTimeout(5 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("new mongodb client error, %s", err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("ping mongodb server(%s) error, %s", m.Endpoints, err)
	}

	return client, nil
}

func newDefaultCache() *_cache {
	return &_cache{
		Type:   "memory",
		Memory: memory.NewDefaultConfig(),
		Redis:  redis.NewDefaultConfig(),
	}
}

type _cache struct {
	Type   string         `toml:"type" json:"type" yaml:"type" env:"CACHE_TYPE"`
	Memory *memory.Config `toml:"memory" json:"memory" yaml:"memory"`
	Redis  *redis.Config  `toml:"redis" json:"redis" yaml:"redis"`
}

func newDefaultEtcd() *Etcd {
	return &Etcd{
		InstanceTTL: 300,
		Prefix:      "inforboard",
	}
}

type Etcd struct {
	Endpoints   []string `toml:"endpoints" env:"ETCD_ENDPOINTS" envSeparator:","`
	UserName    string   `toml:"username" env:"ETCD_USERNAME"`
	Password    string   `toml:"password" env:"ETCD_PASSWORD"`
	Prefix      string   `toml:"prefix" env:"ETCD_Prefix"`
	InstanceTTL int64    `toml:"instance_ttl" env:"ETCD_INSTANCE_TTL"`
}

func (e *Etcd) Validate() error {
	if len(e.Endpoints) == 0 {
		return fmt.Errorf("etcd enpoints not config")
	}
	return nil
}

func (e *Etcd) GetClient() *clientv3.Client {
	if etcdClient == nil {
		panic("please load etcd client first")
	}

	return etcdClient
}

func (e *Etcd) getClient() (*clientv3.Client, error) {
	timeout := time.Duration(5) * time.Second
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: timeout,
		Username:    e.UserName,
		Password:    e.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("connect etcd error, %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ml, err := client.MemberList(ctx)
	if err != nil {
		return nil, err
	}
	zap.L().Debugf("etcd members: %s", ml.Members)

	return client, nil
}
