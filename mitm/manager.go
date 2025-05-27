package mitm

import (
	"fmt"
	"time"
)

// MitmManager 提供高级的mitm管理接口
type MitmManager struct {
	proxyManager *ProxyManager
}

// NewMitmManager 创建新的mitm管理器
func NewMitmManager() *MitmManager {
	return &MitmManager{
		proxyManager: NewProxyManager(),
	}
}

// StartCapture 开始数据捕获
func (mm *MitmManager) StartCapture(processors []MitmProcessor) error {
	if mm.proxyManager.IsRunning() {
		return fmt.Errorf("捕获服务已在运行中")
	}

	return mm.proxyManager.Start(processors)
}

// StopCapture 停止数据捕获
func (mm *MitmManager) StopCapture() error {
	return mm.proxyManager.Stop()
}

// CaptureLogin 专门用于捕获登录信息
func (mm *MitmManager) CaptureLogin(timeout time.Duration) (*LoginInfo, error) {
	if mm.proxyManager.IsRunning() {
		return nil, fmt.Errorf("代理服务已在运行中，请先停止")
	}

	processors := []MitmProcessor{
		NewLoginProcessor(),
	}

	if err := mm.proxyManager.Start(processors); err != nil {
		return nil, fmt.Errorf("启动代理失败: %v", err)
	}

	// 确保退出时清理资源
	defer func() {
		if err := mm.proxyManager.Stop(); err != nil {
			fmt.Printf("停止代理时出错: %v\n", err)
		}
	}()

	// 设置超时
	timeoutChan := time.After(timeout)

	// 等待结果、停止信号或超时
	select {
	case result := <-mm.proxyManager.GetResultChan():
		if loginInfo, ok := result.(*LoginInfo); ok {
			return loginInfo, nil
		}
		return nil, fmt.Errorf("接收到非预期的结果类型")
	case <-mm.proxyManager.GetStopChan():
		return nil, fmt.Errorf("捕获被中断")
	case <-timeoutChan:
		return nil, fmt.Errorf("捕获超时")
	}
}

// GetResults 获取捕获结果通道
func (mm *MitmManager) GetResults() <-chan interface{} {
	return mm.proxyManager.GetResultChan()
}

// IsRunning 检查是否正在运行
func (mm *MitmManager) IsRunning() bool {
	return mm.proxyManager.IsRunning()
}

// 全局管理器实例
var DefaultManager = NewMitmManager()
