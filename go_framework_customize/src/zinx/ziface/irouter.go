package ziface

/*
	路由是干嘛的？提供一个指令，然后这个指令所对应的处理方式，路由和处理方式放在一起叫路由
	路由抽象接口
	路由里的数据都是IRequest
*/

type IRouter interface {
	//在处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)
	//在处理conn业务的主方法
	Handle(request IRequest)
	//在处理conn业务之后的钩子方法Hook
	PostHandle(request IRequest)
}