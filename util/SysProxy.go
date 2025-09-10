package util

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

var (
	originProxyServer string
	originProxyEnable uint32
)

func EnableSysProxy(port int) error {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil && !errors.Is(err, registry.ErrNotExist) {
		return err
	}
	originProxyServer = proxyServer

	if err := key.SetStringValue("ProxyServer", fmt.Sprintf("127.0.0.1:%d", port)); err != nil {
		return err
	}

	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		return err
	}
	originProxyEnable = uint32(proxyEnable)

	if err := key.SetDWordValue("ProxyEnable", 1); err != nil {
		return err
	}

	return nil
}

func DisableSysProxy() error {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Software\Microsoft\Windows\CurrentVersion\Internet Settings`,
		registry.ALL_ACCESS,
	)
	if err != nil {
		return err
	}
	defer key.Close()

	if err := key.SetStringValue("ProxyServer", originProxyServer); err != nil {
		return err
	}

	if err := key.SetDWordValue("ProxyEnable", originProxyEnable); err != nil {
		return err
	}

	return nil
}
