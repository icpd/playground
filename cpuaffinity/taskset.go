package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	pid := os.Getpid()

	// 获取当前进程的 CPU 亲和性
	cmd := exec.Command("taskset", "-p", fmt.Sprintf("%d", pid))
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	// 设置当前进程的 CPU 亲和性为 CPU 0 和 CPU 1
	cmd = exec.Command("taskset", "-p", "0,1", fmt.Sprintf("%d", pid))
	out, err = cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", out)
}
