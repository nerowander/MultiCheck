package main

import (
	"FinalProject/cmd"
)

func main() {
	cmd.Execute()
}

// Dockerfile判断两种界面模式 -> 选择是否模块容器化 -> 指定容器目标和参数（容器需要提前部署）-> 每一个模块配一个main.go和dockerfile
// 一共三个模块，3个http路由，试一下http模式看看后续是否可用，一些默认config参数可通过k8s yaml传参
// 有个小问题是：每一个容器理论上会保存一部分的日志放到文件里，那么需要统一日志名，可以在 Deployment 或 Pod 配置中，使用 HostPath 卷将容器内的某个目录挂载到宿主机的目录
// 3个模块容器，还是1个容器里面有3个模块：可先尝试3个模块容器
