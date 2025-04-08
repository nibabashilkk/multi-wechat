package service

import (
	"bufio"
	"fmt"
	"golang.org/x/sys/windows/registry"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type WXService struct{}

func (*WXService) RunMultiWechat() {
	fmt.Println("=== 微信多开工具 ===")

	// 从注册表获取微信安装路径
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Tencent\WeChat`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Println("错误：无法打开注册表")
		fmt.Println("按回车键退出...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	defer k.Close()

	wechatPath, _, err := k.GetStringValue("InstallPath")
	if err != nil {
		fmt.Println("错误：没有找到微信安装路径")
		fmt.Println("按回车键退出...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	// 获取用户输入
	fmt.Print("请输入要开启的微信数量: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	num, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || num <= 0 {
		fmt.Println("错误：请输入有效的数字")
		fmt.Println("按回车键退出...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	// 结束现有的微信进程
	exec.Command("taskkill", "/f", "/t", "/im", "WeChat.exe").Run()

	// 切换到微信目录
	err = os.Chdir(wechatPath)
	if err != nil {
		fmt.Println("错误：无法切换到微信目录")
		fmt.Println("按回车键退出...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	// 启动多个微信实例
	wechatExe := filepath.Join(wechatPath, "WeChat.exe")
	fmt.Printf("正在启动 %d 个微信实例...\n", num)
	for i := 0; i < num; i++ {
		cmd := exec.Command(wechatExe)
		err = cmd.Start()
		if err != nil {
			fmt.Printf("错误：启动第 %d 个微信实例失败\n", i+1)
		} else {
			fmt.Printf("已启动第 %d 个微信实例\n", i+1)
		}
	}

	fmt.Println("\n所有实例已启动完成")
	fmt.Println("按回车键退出...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
