FROM golang:latest
MAINTAINER Hanya
# 修改国内源
#RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
#RUN sed -i 's/security.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
# 执行命令
#RUN apt-get update --fixing-missing
#RUN apt-get install gcc libc6-dev git lrzsz -y
#RUN curl -sSL https://get.daocloud.io/docker | sh
# 下载项目
RUN go get -u github.com/hanya070603/gluster_exporter_docker
# 定制工作目录
WORKDIR /go/bin
# 对外端口
EXPOSE 9189
# 运行gluster_exporter
ENTRYPOINT ["./gluster_exporter_docker --profile"]
