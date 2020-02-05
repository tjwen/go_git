package ziface

type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Serve()

	//路由功能：给当前的服务器注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32,router IRouter)

	//获取当前连接
	GetConnMgr() IConnManager

	//注册（设置）OnConnStart 钩子函数的方法
	SetOnConnStart(func(connection IConnection))
	//注册（设置）OnConnStop 钩子函数的方法
	SetOnConnStop(func(connection IConnection))
	//调用 OnConnStart钩子函数的方法
	CallOnConnStart(IConnection)
	//调用 OnConnStop钩子函数的方法
	CallOnConnStop(IConnection)
}