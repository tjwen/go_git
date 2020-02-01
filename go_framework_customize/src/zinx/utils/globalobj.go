package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数是可以通过Zinx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer //当前全局的server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		Zinx
	*/
	Version        string //当前Zinx的版本号
	MaxConn        int    //当前服务器主机允许的最大链接数
	MaxPackageSize uint32 //当前Zinx框架数据包的最大值

}

/*
	定义一个全局的对外GlobalObj
*/

var GlobalObject *GlobalObj

/*
	从zinx.json去加载用于自定义的参数
*/
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	//将json文件中数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}

/*
	提供一个init方法，初始化GlobalObject
*/
func init() {
	//如果配置文件没有加载，默认值
	GlobalObject = &GlobalObj{
		TcpServer:      nil,
		Host:           "127.0.0.1",
		TcpPort:        8999,
		Name:           "ZinxServerApp",
		Version:        "V0.6",
		MaxConn:        1000,
		MaxPackageSize: 4096,
	}
	//应该从conf/zinx.json 加载一些用户自定义的参数
	//GlobalObject.Reload()
}
