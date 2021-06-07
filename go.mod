module github.com/infraboard/workflow

go 1.16

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/caarlos0/env/v6 v6.6.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/infraboard/keyauth v0.3.1
	github.com/infraboard/mcube v1.0.5
	github.com/infraboard/protoc-gen-go-ext v0.0.2
	github.com/spf13/cobra v1.1.3
	go.etcd.io/etcd v3.3.25+incompatible
	go.etcd.io/etcd/api/v3 v3.5.0-beta.3
	go.etcd.io/etcd/client/pkg/v3 v3.5.0-beta.3 // indirect
	go.mongodb.org/mongo-driver v1.5.2
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	k8s.io/client-go v0.20.4
)
