# 使用官方 Ubuntu 镜像作为基础镜像
FROM ubuntu:20.04

# 设置环境变量，避免交互式配置
ENV DEBIAN_FRONTEND=noninteractive

# 更新 apt 包索引并安装 openssh-server
RUN apt-get update && \
    apt-get install -y openssh-server sudo && \
    rm -rf /var/lib/apt/lists/*

# 设置 root 用户密码为 admin123
RUN echo 'root:admin123' | chpasswd

# 创建 SSH 目录
RUN mkdir /var/run/sshd

# 允许 root 登录
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

# 设置容器启动时运行 sshd 服务
CMD ["/usr/sbin/sshd", "-D"]