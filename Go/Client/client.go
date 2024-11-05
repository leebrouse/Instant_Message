package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

// Client 结构体，表示一个客户端
type Client struct {
	ServerIp   string   // 服务器IP地址
	ServerPort int      // 服务器端口
	Name       string   // 用户名
	conn       net.Conn // 与服务器的连接
	flag       int      // 用于选择操作的标志
}

// NewClient 创建并初始化一个新的客户端实例，连接到指定的服务器
func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999, // 初始标志值
	}

	// 尝试连接到服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		// 如果连接失败，输出错误信息并返回 nil
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn // 赋值连接对象
	return client
}

// DealRespone 处理服务器的响应，并将其输出到标准输出
func (client *Client) DealRespone() {
	// 将服务器的响应内容复制到标准输出（终端）
	io.Copy(os.Stdout, client.conn)
}

// menu 显示客户端的操作菜单，供用户选择
func (client *Client) menu() bool {
	var flag int

	// 打印菜单
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	// 读取用户输入的选项
	fmt.Scanln(&flag)

	// 检查用户输入是否有效
	if flag >= 0 && flag <= 3 {
		client.flag = flag // 设置操作标志
		return true
	} else {
		// 输入无效，提示重新输入
		fmt.Println(">>>>>>>请输入合法范围的数字<<<<<<")
		return false
	}
}

// SelectUsers 发送 "who" 请求，获取在线用户列表
func (client *Client) SelectUsers() {
	sendMsg := "who\n" // 请求在线用户
	_, err := client.conn.Write([]byte(sendMsg)) // 发送请求
	if err != nil {
		// 如果发送失败，输出错误信息
		fmt.Println("conn Write err:", err)
		return
	}
}

// PrivateChat 进行私聊
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	// 请求在线用户列表
	client.SelectUsers()
	fmt.Println(">>>>>请输入聊天对象[用户名]，exit退出;")
	fmt.Scanln(&remoteName)

	// 如果用户输入不是 "exit" 则继续进行聊天
	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容，exit退出")
		fmt.Scanln(&chatMsg)

		// 循环直到用户输入 "exit" 退出聊天
		for chatMsg != "exit" {
			// 如果消息不为空，发送消息
			if len(chatMsg) != 0 {
				sendMsg := chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg)) // 发送消息
				if err != nil {
					// 发送失败，输出错误信息
					fmt.Println("conn Write err:", err)
					break
				}
			}

			// 清空消息内容，等待下一条消息
			chatMsg = ""
			fmt.Println(">>>>请输入消息内容，exit退出")
			fmt.Scanln(&chatMsg)
		}

		// 继续请求在线用户列表
		client.SelectUsers()
		fmt.Println(">>>>>请输入聊天对象[用户名]，exit退出;")
		fmt.Scanln(&remoteName)
	}
}

// PublicChat 进行公聊
func (client *Client) PublicChat() {
	var chatMsg string

	// 提示用户输入消息内容
	fmt.Println(">>>>请输入聊天内容,exit退出.")
	fmt.Scanln(&chatMsg)

	// 循环直到用户输入 "exit"
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			// 发送消息
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg)) // 发送消息
			if err != nil {
				// 发送失败，输出错误信息
				fmt.Println("conn Write err:", err)
				break
			}
		}

		// 清空消息内容，等待下一条消息
		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容,exit退出.")
		fmt.Scanln(&chatMsg)
	}
}

// UpdateName 更新客户端的用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>>请输入用户名：")
	fmt.Scanln(&client.Name)

	// 构造发送的消息
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg)) // 发送更新用户名的请求
	if err != nil {
		// 如果发送失败，输出错误信息
		fmt.Println("conn.Write err", err)
		return false
	}

	return true
}

// Run 运行客户端程序，根据选择的操作执行相应的功能
func (client *Client) Run() {
	// 当选择不是退出时，继续运行
	for client.flag != 0 {
		// 显示菜单并等待用户输入
		for client.menu() != true {
		}

		// 根据用户的选择，执行相应的操作
		switch client.flag {
		case 1:
			// 公聊模式
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// init 函数，用于解析命令行参数
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认为127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口号(默认为8888)")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	// 创建客户端实例并尝试连接服务器
	client := NewClient(serverIp, serverPort)
	if client == nil {
		// 连接失败，输出错误信息
		fmt.Println(">>>>>>>>链接失败")
		return
	}

	// 启动一个goroutine，处理服务器的响应
	go client.DealRespone()

	// 输出连接成功信息
	fmt.Println(">>>>>>链接成功")

	// 运行客户端程序
	client.Run()
}
