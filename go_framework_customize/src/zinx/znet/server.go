package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

//iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的IP版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int

	//当前的Server添加一个router,server注册的链接对用的处理业务
	Router ziface.IRouter
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[start]Server Listenner at IP :%s, Port:%d , is starting\n", s.IP, s.Port)
	go func() {
		//1创建一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//2监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err:", err)
			return
		}
		fmt.Println("start Zinx server success,", s.Name, "success,Listening...")

		var cid uint32 = 0

		//3阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			//如果有客户链接过来，阻塞会返回

			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:", err)
				continue
			}
		/*	//已经与客户端建立连接，做一些业务
			go func() {
				for {
					buf := make([]byte, 512)
					n, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err:", err)
					}
					//回显
					if _, err := conn.Write(buf[:n]); err != nil {
						fmt.Println("write back buf err:", err)
						continue
					}
				}
			}()*/
			dealConn := NewConnection(conn, cid, s.Router)
			cid ++
			go dealConn.Start()
		}
	}()

}

//停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态、已经开辟的链接信息 进行停止或者回收

}

//运行服务器
func (s *Server) Serve() {
	s.Start()

	//TODO 做一些启动服务之后的额外业务 __小坑

	//阻塞状态
	select {}
}

//路由功能：给当前的服务器注册一个路由方法，供客户端的链接处理使用
func(s *Server) AddRouter(router ziface.IRouter){
	s.Router = router
	fmt.Println("Add Router Success!")
}

/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {

	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
