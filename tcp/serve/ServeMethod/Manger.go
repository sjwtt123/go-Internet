package ServeMethod

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

var HbManager = NewHeartbeatManager(30*time.Second, 10*time.Second) //创建加锁客户端实例

type HeartbeatManager struct {
	clients  map[string]*Client // 存储所有客户端
	mutex    sync.RWMutex       // 读写锁，保证并发安全
	timeout  time.Duration      // 心跳超时时间
	interval time.Duration      // 检查间隔
}

// NewHeartbeatManager 创建心跳管理器
func NewHeartbeatManager(timeout, interval time.Duration) *HeartbeatManager {
	return &HeartbeatManager{
		clients:  make(map[string]*Client),
		timeout:  timeout,
		interval: interval,
	}
}

// AddClient 添加客户端
func (hm *HeartbeatManager) AddClient(client *Client) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	client.LastActive = time.Now() // 设置初始活跃时间
	hm.clients[client.Nickname] = client
	fmt.Printf("客户端 %s 已添加到心跳检测\n", client.Nickname)

}

// RemoveClient 移除客户端
func (hm *HeartbeatManager) RemoveClient(client *Client) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	delete(hm.clients, client.Nickname)
	fmt.Printf("客户端 %s 已从心跳检测移除\n", client.Nickname)
}

// UpdateClientActivity 更新客户端活跃时间
func (hm *HeartbeatManager) UpdateClientActivity(name string) {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()
	if HbManager.clients[name] == nil {
		return
	}

	HbManager.clients[name].LastActive = time.Now()

}

// Start 启动心跳检测 每10秒触发一次
func (hm *HeartbeatManager) Start() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Start()协程发生 panic: %v\n", r)
			debug.PrintStack() // 打印堆栈跟踪
		}
	}()

	ticker := time.NewTicker(hm.interval)
	defer ticker.Stop()

	fmt.Printf("心跳检测启动 - 超时: %v, 检查间隔: %v\n", hm.timeout, hm.interval)

	for range ticker.C {
		hm.checkHeartbeats()
	}
}

// checkHeartbeats 检查所有客户端的心跳
func (hm *HeartbeatManager) checkHeartbeats() {
	hm.mutex.Lock()
	defer hm.mutex.Unlock()

	now := time.Now()

	for _, client := range hm.clients {
		// 计算距离最后活跃时间的时间差
		timeSinceLastActive := now.Sub(client.LastActive)

		// 如果超时，关闭连接并移除
		if timeSinceLastActive > hm.timeout {
			log.Printf("客户端 %s 心跳超时 (%.0f秒) 断开连接\n",
				client.Nickname, timeSinceLastActive.Seconds())
			leaveChan <- client // 从管理器中移除
		} else {
			// 正常情况，打印心跳状态
			fmt.Printf("客户端 %s 心跳正常 (%.0f秒前活跃)\n",
				client.Nickname, timeSinceLastActive.Seconds())
		}
	}
}
