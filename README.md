shipper 原仓库地址为：https://github.com/EwanValentine/shippy。

shipper 的技术栈包括 gRPC、protobuf、docker、go-micro 等，是一个不错的微服务学习项目。作者编写了一个系列的文章，循序渐进的展示了项目搭建的过程，系列文章的地址 [点击这里](https://ewanvalentine.io/microservices-in-golang-part-1)。

此项目用于学习，然后根据自身的理解添加合适的中文注释，原系列有部分中文翻译的文章，但并不完善而且内容与英语原文稍有差异（可能是原文有更新，翻译的文章没有跟进）。

注意：

tutorial-2 在 docker 容器中执行 `go module download` 命令，会因为网络问题导致 golang.org 的包下载失败，我是依照 [此文章](https://www.simpleapples.com/2019/04/18/building-docker-image-behind-proxy/) 解决的代理问题。

每个阶段都对应一个分支 branch：

### tutorial-1
基于 gRPC 和 protobuf 完成一个最简单的微服务。
### tutorial-2
在 tutorial-1 的基础上，使用 go-micro 和 docker 服务化，并添加了 vessel 货船服务。

货运服务中使用了vessel的服务，在 docker 中 go build 会出现找不到 vessel 包的情况，需要先将 vessel 上传到 github，然后修改导入地址为线上才可以。

### 名词注释

每次货运服务需要将集装箱放到合适的货船上。

- consignment: 货运服务
- vessel: 货船
- container: 集装箱