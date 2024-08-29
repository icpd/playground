package main

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	var mask uintptr

	// 获取当前进程的 CPU 亲和性
	if _, _, err := syscall.RawSyscall(syscall.SYS_SCHED_GETAFFINITY, 0, uintptr(unsafe.Sizeof(mask)), uintptr(unsafe.Pointer(&mask))); err != 0 {
		fmt.Println("failed to get CPU affinity:", err)
		return
	}
	fmt.Println("current CPU affinity:", mask)

	// 设置当前进程的 CPU 亲和性为 CPU 0 和 CPU 1
	mask = 3
	if _, _, err := syscall.RawSyscall(syscall.SYS_SCHED_SETAFFINITY, 0, uintptr(unsafe.Sizeof(mask)), uintptr(unsafe.Pointer(&mask))); err != 0 {
		fmt.Println("failed to set CPU affinity:", err)
		return
	}
	fmt.Println("new CPU affinity:", mask)

	for {
		println("Hello, World!")
		time.Sleep(1 * time.Second)
	}
}
