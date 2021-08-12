module github.com/infraboard/workflow

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/caarlos0/env/v6 v6.6.0
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/go-playground/validator/v10 v10.9.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/infraboard/keyauth v0.5.4
	github.com/infraboard/mcube v1.3.7
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/rs/xid v1.3.0
	github.com/spf13/cobra v1.2.1
	github.com/stretchr/testify v1.7.0
	go.etcd.io/etcd/api/v3 v3.5.0
	go.etcd.io/etcd/client/v3 v3.5.0
	go.mongodb.org/mongo-driver v1.7.1
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.27.1
	k8s.io/apimachinery v0.20.6
	k8s.io/client-go v0.20.6
)
