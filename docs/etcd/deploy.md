# 部署


## 单机部署

```sh
mkdir /data/etcd -p
docker run -d \
  -p 32379:2379 \
  -p 32380:2380 \
  -v /data/etcd:/etcd-data/member \
  --name exam-etcd \
   quay.io/coreos/etcd:latest \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new
```

## 常用操作
```sh
# 添加
docker exec -e ETCDCTL_API=3 exam-etcd etcdctl --endpoints=http://127.0.0.1:2379 put foo bar
# 查看
docker exec -e ETCDCTL_API=3 exam-etcd etcdctl --endpoints=http://127.0.0.1:2379  get --prefix foo
# 删除
docker exec -e ETCDCTL_API=3 exam-etcd etcdctl --endpoints=http://127.0.0.1:2379  del foo
```