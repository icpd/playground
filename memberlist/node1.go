package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/memberlist"
)

type MyEventDelegate1 struct {
	nodeName string
}

func (d *MyEventDelegate1) NotifyJoin(node *memberlist.Node) {
	log.Printf("[%s] *** 通知: 节点加入了: %s (%s:%d)\n", d.nodeName, node.Name, node.Addr, node.Port)
}

func (d *MyEventDelegate1) NotifyLeave(node *memberlist.Node) {
	log.Printf("[%s] *** 通知: 节点离开了 (或被标记为故障): %s (%s:%d)\n", d.nodeName, node.Name, node.Addr, node.Port)
	// 在这里可以执行对故障节点的处理逻辑，例如从服务发现中移除
}

func (d *MyEventDelegate1) NotifyUpdate(node *memberlist.Node) {
	log.Printf("[%s] *** 通知: 节点元数据更新了: %s (元数据: %s)\n", d.nodeName, node.Name, string(node.Meta))
}

func main() {
	nodeName := "node-A"
	bindPort := 7946 // Node A 的监听端口

	// 1. 配置 memberlist
	config := memberlist.DefaultLocalConfig()
	config.Name = nodeName
	config.BindAddr = "0.0.0.0" // 监听所有可用网络接口
	config.BindPort = bindPort
	config.LogOutput = os.Stdout                          // 输出 memberlist 内部日志
	config.Events = &MyEventDelegate1{nodeName: nodeName} // 注册事件委托

	// 为了更快的检测，我们可以缩短一些超时（生产环境不建议过于激进）
	config.ProbeInterval = 1 * time.Second
	config.ProbeTimeout = 500 * time.Millisecond
	config.SuspicionMult = 1 // 降低怀疑乘数，更快进入 Dead 状态

	// 2. 创建 memberlist 实例
	list, err := memberlist.Create(config)
	if err != nil {
		log.Fatalf("创建 memberlist 失败: %v", err)
	}
	defer list.Shutdown() // 确保程序退出时优雅关闭

	log.Printf("[%s] 节点已启动，监听在 %s:%d\n", nodeName, config.BindAddr, bindPort)
	log.Printf("[%s] 等待其他节点加入...\n", nodeName)

	// 3. 打印当前成员列表 (定期)
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			fmt.Printf("\n[%s] 当前集群成员 (%d个):\n", nodeName, list.NumMembers())
			for _, member := range list.Members() {
				// 获取成员的健康状态
				status := "Alive"
				if member.State == memberlist.StateSuspect {
					status = "Suspect"
				} else if member.State == memberlist.StateDead {
					status = "Dead"
				}
				fmt.Printf("  - 名称: %s, 地址: %s:%d, 状态: %s\n", member.Name, member.Addr, member.Port, status)
			}
			fmt.Println("-------------------------------------------")
		}
	}()

	// 4. 等待信号，以便优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Printf("[%s] 收到退出信号，正在关闭 memberlist...\n", nodeName)
}
