package main

import (
	"log"
	"os"
	"time"

	"github.com/hashicorp/memberlist"
)

type MyEventDelegate3 struct {
	nodeName string
}

func (d *MyEventDelegate3) NotifyJoin(node *memberlist.Node) {
	log.Printf("[%s] *** 通知: 节点加入了: %s (%s:%d)\n", d.nodeName, node.Name, node.Addr, node.Port)
}

func (d *MyEventDelegate3) NotifyLeave(node *memberlist.Node) {
	log.Printf("[%s] *** 通知: 节点离开了 (或被标记为故障): %s (%s:%d)\n", d.nodeName, node.Name, node.Addr, node.Port)
}

func (d *MyEventDelegate3) NotifyUpdate(node *memberlist.Node) {
	log.Printf("[%s] *** 通知: 节点元数据更新了: %s (元数据: %s)\n", d.nodeName, node.Name, string(node.Meta))
}

func main() {
	nodeName := "node-C"
	bindPort := 7948             // Node C 的监听端口
	joinAddr := "127.0.0.1:7946" // Node A 的地址

	// 1. 配置 memberlist
	config := memberlist.DefaultLocalConfig()
	config.Name = nodeName
	config.BindAddr = "0.0.0.0"
	config.BindPort = bindPort
	config.LogOutput = os.Stdout
	config.Events = &MyEventDelegate3{nodeName: nodeName}

	// 和 Node A 保持一致的探测配置
	config.ProbeInterval = 1 * time.Second
	config.ProbeTimeout = 500 * time.Millisecond
	config.SuspicionMult = 1

	// 2. 创建 memberlist 实例
	list, err := memberlist.Create(config)
	if err != nil {
		log.Fatalf("创建 memberlist 失败: %v", err)
	}
	defer list.Shutdown() // 确保程序退出时优雅关闭

	log.Printf("[%s] 节点已启动，监听在 %s:%d\n", nodeName, config.BindAddr, bindPort)

	// 3. 加入集群
	log.Printf("[%s] 尝试加入集群 (通过 %s)...\n", nodeName, joinAddr)
	_, err = list.Join([]string{joinAddr})
	if err != nil {
		log.Fatalf("加入集群失败: %v", err)
	}
	log.Printf("[%s] 成功加入集群。\n", nodeName)
	select {}

	// 4. 模拟运行一段时间后自动关闭
	runDuration := 15 * time.Second // 运行 15 秒后自动关闭
	log.Printf("[%s] 将运行 %s 后自动下线，以模拟故障或优雅退出。\n", nodeName, runDuration)

	// 模拟人工关闭或程序崩溃，这里用 Sleep 和 Shutdown 模拟
	time.Sleep(runDuration)

	log.Printf("[%s] 模拟下线：主动调用 list.Shutdown()。\n", nodeName)
	// list.Shutdown() 会向集群发送 Leave 消息，然后关闭。
	// 如果不调用 Shutdown，直接 os.Exit()，memberlist 将会通过故障检测机制来发现其失效。
	// 这里使用 Shutdown 模拟“优雅退出”，Node A 也会收到 NotifyLeave。
	// 如果要模拟“崩溃”，可以直接在这里 os.Exit(1)
	os.Exit(1)
	list.Shutdown()

	log.Printf("[%s] 节点已停止运行。\n", nodeName)

	// 保持程序活着一段时间，确保 Shutdown 消息有时间发送出去
	time.Sleep(2 * time.Second)
}
