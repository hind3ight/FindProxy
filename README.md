# FindProxy
A simple demo for finding proxy.

## 代码开发计划v1
1. 读取文件（大陆网段）
2. 根据网段获取所有ip (共4,213,413,773,312个端口需要探测)
3. 根据ip地址，扫描其可用端口，并存储
4. 将ip地址加端口进行代理测试，并持久化

## 代码开发计划v2
1. 根据可用的地址分别对所需服务进行代理测试
2. 代理地址与服务地址做关联，并打上优先级标签