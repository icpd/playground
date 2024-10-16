/*
index = head & (len(queue) - 1)
这个位运算的作用是让 head 永远保持在 [0, len(queue) - 1] 的范围内，而不用复杂的模运算 %。
因为队列的长度是 2 的幂（如 8、16），所以 len(queue) - 1 会变成全 1 的二进制数（如 7 或 15），这样 & 操作就等价于模运算。
https://mp.weixin.qq.com/s/fj87oGZPkFKQiGZxhrYRVQ
*/
package main

import "fmt"

type RingBuffer struct {
	buffer []int
	head   int
	tail   int
	size   int
}

func NewRingBuffer(size int) *RingBuffer {
	// 是否为 2 的幂
	if size > 0 && (size&(size-1)) != 0 {
		panic("size 必须是 2 的幂")
	}

	return &RingBuffer{
		buffer: make([]int, size),
		size:   size,
	}
}

func (rb *RingBuffer) PushHead(val int) {
	rb.buffer[rb.head&(rb.size-1)] = val
	rb.head++
}

func (rb *RingBuffer) PopTail() int {
	val := rb.buffer[rb.tail&(rb.size-1)]
	rb.tail++
	return val
}

func main() {
	rb := NewRingBuffer(8) // 大小为 8 的环形缓冲区
	rb.PushHead(1)
	rb.PushHead(2)
	rb.PushHead(3)

	fmt.Println(rb.PopTail()) // 输出 1
	fmt.Println(rb.PopTail()) // 输出 2
}
