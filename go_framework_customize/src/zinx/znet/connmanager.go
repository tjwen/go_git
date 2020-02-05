package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	Connections map[uint32]ziface.IConnection //管理的连接集合
	connLock    sync.RWMutex                  //保护连接集合的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: make(map[uint32]ziface.IConnection),
	}
}

//添加链接
func (cm *ConnManager) Add(conn ziface.IConnection) {
	//保护共享资源map,加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//将conn加入到ConnManager中
	cm.Connections[conn.GetConnID()] = conn
	fmt.Println("connID= ",conn.GetConnID(),"connection add to ConnManager successfully: conn num = ", cm.Len())
}

//删除链接
func (cm *ConnManager) Remove(conn ziface.IConnection) {
	//保护共享资源map,加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//删除
	delete(cm.Connections,conn.GetConnID())
	fmt.Println("connID= ",conn.GetConnID(),"connection remove from ConnManager successfully: conn num = ", cm.Len())
}

//根据ID获取链接
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	//保护共享资源map,加读锁
	cm.connLock.RLock()
	defer cm.connLock.Unlock()
	if conn,ok:= cm.Connections[connID]; ok {
		return  conn,nil
	}else{
		return nil,errors.New("connection not FOUND")
	}
}

//得到当前的链接数
func (cm *ConnManager) Len() int {
	return len(cm.Connections)
}

//清除并终止所有的d连接
func (cm *ConnManager) ClearConn() {
	//保护共享资源map,加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//删除conn并停止 conn的工作
	for connID,conn := range cm.Connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.Connections,connID)
	}
	fmt.Println("Clear All connections success")
}
