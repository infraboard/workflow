# workflow

应用研发交付中心, 以应用为中心的分布式自定义流水线, 支持的常见场景包括:
+ 应用研发流, 常见CI CD
+ 应用审批流, 比如应用资源审批与自动化流程处理编排
+ 应用事件流, 基于应用事件, 为其编排应急预案流水线, 比如扩容预案流水线和缩容预案流水线


## 架构图


## 快速开发

## 注意

依赖的 mod包: go.etcd.io/etcd v3.3.25+incompatible, 直接替换的v3.5.0-beta.3

make脚手架
```sh
➜  workflow git:(master) ✗ make help
dep                            Get the dependencies
lint                           Lint Golang files
vet                            Run go vet
test                           Run unittests
test-coverage                  Run tests with coverage
build                          Local build
linux                          Linux build
run                            Run Server
clean                          Remove previous build
help                           Display this help screen
```

1. 使用go mod下载项目依赖
```sh
$ make dep
```

2. 添加配置文件(默认读取位置: etc/workflow.toml)
```sh
$ 编辑样例配置文件 etc/workflow.toml.example
$ mv etc/workflow.toml.example etc/workflow.toml
```

3. 启动服务
```sh
$ make run
```

## 相关文档