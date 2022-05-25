# 应用部署


## 注入变量管理

```yaml
    env:
    - name: POD_OWN_IP_ADDRESS
      valueFrom:
        fieldRef:
          fieldPath: status.podIP
    - name: POD_OWN_NAME
      valueFrom:
        fieldRef:
          fieldPath: metadata.name
    - name: POD_OWN_NAMESPACE
      valueFrom:
        fieldRef:
          fieldPath: metadata.namespace
```


```
APP_NODE_NAME #部署节点名字 [spec.nodeName]      
APP_POD_NAME #部署POD名字 [metadata.name]
APP_POD_NAMESPACE #应用命名空间 [metadata.namespace]
APP_POD_IP #部署POD_IP地址信息 [status.podIP]
APP_POD_SERVICE_ACCOUNT #部署服务用户名称 [spec.serviceAccountName]
```
