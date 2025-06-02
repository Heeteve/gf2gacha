package util

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
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

	if err := key.SetStringValue("ProxyServer", fmt.Sprintf("127.0.0.1:%d", port)); err != nil {
		return err
	}

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

	if err := key.SetDWordValue("ProxyEnable", 0); err != nil {
		return err
	}

	return nil
}
