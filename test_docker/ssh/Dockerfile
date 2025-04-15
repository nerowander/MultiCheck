# 使用官方 Ubuntu 基础镜像
FROM ubuntu:20.04

# 安装 OpenSSH 服务

#RUN apt-get update && \
#    apt-get -y install ca-certificates && \
#    mv /etc/apt/sources.list /etc/apt/sources_backup.list && \
#    echo "deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ focal main restricted universe multiverse" >> /etc/apt/sources.list && \
#    echo "deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ focal-updates main restricted universe multiverse" >> /etc/apt/sources.list && \
#    echo "deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list && \
#    echo "deb http://security.ubuntu.com/ubuntu/ focal-security main restricted universe multiverse" >> /etc/apt/sources.list && \
##    echo "deb https://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse" >> /etc/apt/sources.list && \
##    echo "deb-src https://mirrors.aliyun.com/ubuntu/ focal-updates main restricted universe multiverse" >> /etc/apt/sources.list && \
##    echo "# deb https://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse" >> /etc/apt/sources.list && \
##    echo "# deb-src https://mirrors.aliyun.com/ubuntu/ focal-proposed main restricted universe multiverse" >> /etc/apt/sources.list && \
##    echo "deb https://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list && \
##    echo "deb-src https://mirrors.aliyun.com/ubuntu/ focal-backports main restricted universe multiverse" >> /etc/apt/sources.list && \
##    echo "deb http://archive.canonical.com/ubuntu focal partner " >> /etc/apt/sources.list && \
#    apt-get update -y && apt-get install -y openssh-server && \
#    echo 'root:admin123' | chpasswd && \
#    sed -i 's/^#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
#    mkdir /var/run/sshd

# 暴露 SSH 端口
EXPOSE 22

# 启动 SSH 服务
#CMD ["/usr/sbin/sshd", "-D"]