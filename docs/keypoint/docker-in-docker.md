# 如何为容器提供Docker工具


简单而言: 把宿主机的Docker工具挂载到容器里面, 扩展而言，我们可以在宿主机层面做一个工具箱, 上层容器按需挂载
```
# docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker dockerindocker:1.0 /bin/bash
```

参考: [docker 嵌套技术](https://blog.csdn.net/shida_csdn/article/details/79812817)