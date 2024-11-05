package main

func main() {
	// 创建一个新的服务器实例，指定IP为"127.0.0.1"和端口为8888
	server := NewServer("127.0.0.1", 8888)

	// 启动服务器
	server.Start()
}
