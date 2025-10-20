package ServeMethod

import (
	"sync"
)

var clientManager = NewClientManager() //创建加锁客户端实例

// ClientManager 互斥锁
type ClientManager struct {
	sync.Mutex
	sync.Map
	clients map[Client]bool // 在线客户端
}

// AddClient 添加进入聊天室的客户实例
func (cm *ClientManager) AddClient(cli Client) {
	cm.Lock()
	defer cm.Unlock()
	cm.clients[cli] = true
}

// RemoveClient 移除在聊天室的客户
func (cm *ClientManager) RemoveClient(cli Client) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.clients, cli)
}

// NewClientManager 初始化加锁在线客户端
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[Client]bool), // 初始化
	}
}
