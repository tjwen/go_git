package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
基于Zinx框架来开发的服务器应用程序
*/

//ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

/*//Test PreRouter
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle ...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err!=nil {
		fmt.Println("callback ping error")
		return
	}
}*/
//Test Handle
func (this *PingRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle ...")
	//先读取客户端的数据，再回写pint...ping...ping
	fmt.Println("recv from client : " +
		"msgID:=",request.GetMsgID(),
		"data=",string(request.GetData()))
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil{
		fmt.Println(err)
	}
}
//Test PostHandle
//func (this *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("Call Router PostHandle ...")
//	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
//	if err!=nil {
//		fmt.Println("callback after ping error")
//		return
//	}
//}
type HelloZinxRouter struct {
	znet.BaseRouter
}
func (this *HelloZinxRouter)Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle ...")
	//先读取客户端的数据，再回写pint...ping...ping
	fmt.Println("recv from client : " +
		"msgID:=",request.GetMsgID(),
		"data=",string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome to Zinx"))
	if err != nil{
		fmt.Println(err)
	}
}

func main() {
	//1创建server句柄，使用Zinx的API
	s := znet.NewServer("[zinx v0.6]")

	//给当前zinx框架添加一个自定义的router
	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloZinxRouter{})
	//2启动server
	s.Serve()
	//

}



