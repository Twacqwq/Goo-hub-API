// Package console 命令行辅助方法
package console

import (
	"fmt"
	"os"

	"github.com/mgutz/ansi"
)

// 打印一条成功信息
func Success(msg string) {
	colorOut(msg, "green")
}

// 打印一条错误信息
func Error(msg string) {
	colorOut(msg, "red")
}

// 打印一条警告信息
func Warning(msg string) {
	colorOut(msg, "yellow")
}

func Exit(msg string) {
	Error(msg)
	os.Exit(1)
}

// ExitIf 语法糖 自带 err != nil 判断
func ExitIf(err error) {
	if err != nil {
		Exit(err.Error())
	}
}

// colorOut 内部使用, 设置高亮颜色
func colorOut(message, color string) {
	fmt.Fprintln(os.Stdout, ansi.Color(message, color))
}
