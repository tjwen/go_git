package ziface

/*
	连接管理模块抽象层
*/

type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//根据ID获取链接
	Get(connID uint32)(IConnection,error)
	//得到当前的链接数
	Len() int
	//清除并终止所有的d连接
	ClearConn()
}