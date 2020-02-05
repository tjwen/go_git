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
	//Router ziface.IRouter
	MsgHandler ziface.IMsgHandle

	//该server的连接管理器
	ConnMgr ziface.IConnManager

	//该server创建链接之后自动调用Hook函数
	OnConnStart func(conn ziface.IConnection)
	//该server销毁链接之前调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
}

//定义当前客户端链接的所绑定的业务 handle api(目前这个handle是写死的，以后优化用户自定义handle方法
/*func CallBackToClient(conn *net.TCPConn, data []byte,cnt int) error{
	if _,err := conn.Write(data); err !=nil {
		fmt.Println("write back buf err :",err)
		return errors.New("CallBackToClient error")
	}
	return nil
}
*/

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[start]Server Listenner at IP :%s, Port:%d , is starting\n", s.IP, s.Port)
	go func() {
		s.MsgHandler.StartWorkerPool()

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
			//设置最大连接个数的判断，如果超过最大连接，则关闭此新的连接
			if s.ConnMgr.Len()>= utils.GlobalObject.MaxConn{
				//TODO 给客户端响应一个超出最大连接的错误包
				fmt.Println("Too Many Connections MaxConn =",utils.GlobalObject.MaxConn)
				conn.Close()
				continue
			}
			dealConn := NewConnection(s,conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()

}

//停止服务器
func (s *Server) Stop() {
	// 将一些服务器的资源、状态、已经开辟的链接信息 进行停止或者回收
	fmt.Println("[stop] Zinx Server name:", s.Name)
	s.ConnMgr.ClearConn()
}

//运行服务器
func (s *Server) Serve() {
	s.Start()

	//TODO 做一些启动服务之后的额外业务 __小坑

	//阻塞状态
	select {}
}

//路由功能：给当前的服务器注册一个路由方法，供客户端的链接处理使用
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!")
}

func (s *Server)GetConnMgr() ziface.IConnManager  {
	return s.ConnMgr
}


/*
初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {

	s := &Server{
		Name:       name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}


//注册（设置）OnConnStart 钩子函数的方法
func (s *Server)SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
//注册（设置）OnConnStop 钩子函数的方法
func (s *Server)SetOnConnStop(hookFunc func(connection ziface.IConnection)){
	s.OnConnStop = hookFunc

}
//调用 OnConnStart钩子函数的方法
func (s *Server)CallOnConnStart(conn ziface.IConnection){
	if s.OnConnStart != nil{
		fmt.Println("------>Call OnConnStart()...")
		s.OnConnStart(conn)
	}
}
//调用 OnConnStop钩子函数的方法
func (s *Server)CallOnConnStop(conn ziface.IConnection){
	if s.OnConnStop != nil {
		fmt.Println("----->Call OnConnStop()...")
		s.OnConnStop(conn)
	}
}
