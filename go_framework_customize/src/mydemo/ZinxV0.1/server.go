package main

import "zinx/znet"

/*
基于Zinx框架来开发的服务器应用程序
*/
func main() {
	//1创建server句柄，使用Zinx的API
	server := znet.NewServer("[zinx v0.1]")
	//2启动server
	server.Serve()


}



