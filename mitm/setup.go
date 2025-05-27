package mitm

import (
	"fmt"
	"net/url"

	"golang.org/x/sys/windows/registry"
)

const internetSettingsPath = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`

// SetSystemProxy 设置Windows系统代理
func SetSystemProxy(proxyURL *url.URL) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, internetSettingsPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %v", err)
	}
	defer k.Close()

	if err := k.SetDWordValue("ProxyEnable", 1); err != nil {
		return fmt.Errorf("启用代理失败: %v", err)
	}
	if err := k.SetStringValue("ProxyServer", proxyURL.Host); err != nil {
		return fmt.Errorf("设置代理地址失败: %v", err)
	}
	return nil
}

// ClearSystemProxy 清理Windows系统代理
func ClearSystemProxy() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, internetSettingsPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %v", err)
	}
	defer k.Close()

	if err := k.SetDWordValue("ProxyEnable", 0); err != nil {
		return fmt.Errorf("禁用代理失败: %v", err)
	}
	if err := k.SetStringValue("ProxyServer", ""); err != nil {
		return fmt.Errorf("清除代理地址失败: %v", err)
	}
	return nil
}
