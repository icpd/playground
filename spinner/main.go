package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func main() {
	chart := []string{"⠏", "⠛", "⠹", "⠼", "⠴", "⠦", "⠧"}
	delay := 100 * time.Millisecond

	s := spinner.New(chart, delay)
	s.FinalMSG = "⠿ docker image Pull complete"

	s.Start()

	for i := 0; i < 100; i++ {
		s.Suffix = fmt.Sprintf(" docker image Pulling %dMiB", i)
		time.Sleep(delay)
	}

	s.Stop()

	fmt.Println()

	// =====================
	//       自实现
	// =====================

	spinnerChars := chart               // 旋转动画字符
	totalProgress := 100                // 假设总进度为100
	fmt.Println("Pulling docker image") // 输出任务描述

	// 隐藏光标
	fmt.Print("\033[?25l")

	for progress := 0; progress <= totalProgress; progress++ {
		// 选择旋转字符，或在结束时用满状态符号
		var spinnerChar string
		if progress == totalProgress {
			spinnerChar = "⠿" // 满状态符号
		} else {
			spinnerChar = spinnerChars[progress%len(spinnerChars)]
		}

		bar := fmt.Sprintf("[%s%s]", getBar(progress), getSpaces(50-progress/2))
		fmt.Printf("\r%s %s %d%%", spinnerChar, bar, progress) // 使用 \r 返回行首刷新内容
		time.Sleep(delay)                                      // 控制更新速度
	}

	fmt.Println("\nPull complete!")

}

// getBar 用于生成进度条已完成部分
func getBar(progress int) string {
	completed := progress / 2
	bar := ""
	for i := 0; i < completed; i++ {
		bar += "="
	}
	return bar
}

// getSpaces 用于生成进度条未完成部分
func getSpaces(count int) string {
	spaces := ""
	for i := 0; i < count; i++ {
		spaces += " "
	}
	return spaces
}
