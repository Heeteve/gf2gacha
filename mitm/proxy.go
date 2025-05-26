package mitm

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/lqqyt2423/go-mitmproxy/proxy"
)

// ProxyManager 代理管理器
type ProxyManager struct {
	proxy      *proxy.Proxy
	mutex      sync.Mutex
	isRunning  bool
	stopChan   chan struct{}
	resultChan chan interface{}
	processors []MitmProcessor
}

// NewProxyManager 创建新的代理管理器
func NewProxyManager() *ProxyManager {
	return &ProxyManager{
		resultChan: make(chan interface{}, 10),
	}
}

// MitmAddon mitm代理的插件
type MitmAddon struct {
	proxy.BaseAddon
	processors []MitmProcessor
	resultChan chan interface{}
}

// NewMitmAddon 创建一个新的MitmAddon实例
func NewMitmAddon(processors []MitmProcessor, resultChan chan interface{}) *MitmAddon {
	return &MitmAddon{
		processors: processors,
		resultChan: resultChan,
	}
}

// Response 处理代理响应
func (a *MitmAddon) Response(f *proxy.Flow) {
	capture := &MitmCapture{
		URL:      f.Request.URL.String(),
		Headers:  make(map[string]string),
		Response: f.Response.Body,
	}

	// 复制响应头
	for k, v := range f.Response.Header {
		if len(v) > 0 {
			capture.Headers[k] = v[0]
		}
	}

	// 遍历所有处理器
	for _, processor := range a.processors {
		if processor.Match(capture.URL) {
			if result, err := processor.Process(capture); err == nil {
				select {
				case a.resultChan <- result:
				default:
					// 通道满时不阻塞
				}
			}
		}
	}
}

// Start 启动代理服务
func (pm *ProxyManager) Start(processors []MitmProcessor) error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if pm.isRunning {
		return fmt.Errorf("代理服务已在运行中")
	}

	pm.processors = processors
	pm.stopChan = make(chan struct{})

	// 创建代理选项
	opts := &proxy.Options{
		Addr:        ":9080",
		Debug:       1,
		SslInsecure: true,
		CaRootPath:  "./cert",
	}

	var err error
	pm.proxy, err = proxy.NewProxy(opts)
	if err != nil {
		return fmt.Errorf("创建代理实例失败: %v", err)
	}

	// 设置 addon
	addon := NewMitmAddon(pm.processors, pm.resultChan)
	pm.proxy.AddAddon(addon)

	// 启动代理服务
	go func() {
		if err := pm.proxy.Start(); err != nil {
			fmt.Printf("代理服务错误: %v\n", err)
		}
	}()

	// 设置系统代理
	proxyURL, _ := url.Parse("http://127.0.0.1:9080")
	if err := SetSystemProxy(proxyURL); err != nil {
		pm.proxy.Close()
		return fmt.Errorf("设置系统代理失败: %v", err)
	}

	pm.isRunning = true
	return nil
}

// Stop 停止代理服务
func (pm *ProxyManager) Stop() error {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	if !pm.isRunning {
		return fmt.Errorf("代理服务未在运行")
	}

	if pm.proxy != nil {
		pm.proxy.Close()
		pm.proxy = nil
	}

	if pm.stopChan != nil {
		close(pm.stopChan)
		pm.stopChan = nil
	}

	// 清理系统代理
	if err := ClearSystemProxy(); err != nil {
		return fmt.Errorf("清理系统代理失败: %v", err)
	}

	pm.isRunning = false
	return nil
}

// GetResultChan 获取结果通道
func (pm *ProxyManager) GetResultChan() <-chan interface{} {
	return pm.resultChan
}

// GetStopChan 获取停止通道
func (pm *ProxyManager) GetStopChan() <-chan struct{} {
	return pm.stopChan
}

// IsRunning 检查是否正在运行
func (pm *ProxyManager) IsRunning() bool {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return pm.isRunning
}

// 全局代理管理器实例
var globalProxyManager = NewProxyManager()

// StartMitmProxy 启动mitm代理 (兼容旧接口)
func StartMitmProxy(processors []MitmProcessor, resultChan chan interface{}) error {
	return globalProxyManager.Start(processors)
}

// StopMitmProxy 停止mitm代理 (兼容旧接口)
func StopMitmProxy() error {
	return globalProxyManager.Stop()
}

// StartMitmProxyForLogin 启动mitm代理专门用于登录抓取
func StartMitmProxyForLogin() (*LoginInfo, error) {
	processors := []MitmProcessor{
		NewLoginProcessor(),
	}

	if err := globalProxyManager.Start(processors); err != nil {
		return nil, err
	}

	// 等待结果或停止信号
	select {
	case result := <-globalProxyManager.GetResultChan():
		if loginInfo, ok := result.(*LoginInfo); ok {
			return loginInfo, nil
		}
		return nil, fmt.Errorf("接收到非预期的结果类型")
	case <-globalProxyManager.GetStopChan():
		return nil, nil
	}
}
