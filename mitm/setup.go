package mitm

import (
"encoding/pem"
"fmt"
"net/url"
"os"
"syscall"
"unsafe"

"golang.org/x/sys/windows"
"golang.org/x/sys/windows/registry"
)

var (
	crypt32                              = syscall.NewLazyDLL("crypt32.dll")
	procCertOpenStore                    = crypt32.NewProc("CertOpenStore")
	procCertCloseStore                   = crypt32.NewProc("CertCloseStore")
	procCertAddEncodedCertificateToStore = crypt32.NewProc("CertAddEncodedCertificateToStore")
)

func openCertStore() (windows.Handle, error) {
	storeName, _ := syscall.UTF16PtrFromString("Root")

	r, _, err := procCertOpenStore.Call(uintptr(10), 0, 0, uintptr(0x00010000), uintptr(unsafe.Pointer(storeName)))
	if r == 0 {
		return 0, fmt.Errorf("调用CertOpenStore出错: %v", err)
	}
	return windows.Handle(r), nil
}

func closeCertStore(store windows.Handle) {
	if store != 0 {
		procCertCloseStore.Call(uintptr(store), 0)
	}
}

// InstallMitmCert 安装mitmproxy证书
func InstallMitmCert() error {
	certPath := "./cert/mitmproxy-ca-cert.cer"
	certPem, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("读取证书文件失败: %v", err)
	}

	certDer, _ := pem.Decode(certPem)
	if certDer == nil {
		return fmt.Errorf("pem解码失败")
	}

	store, err := openCertStore()
	if err != nil {
		return err
	}
	defer closeCertStore(store)

	r, _, err := procCertAddEncodedCertificateToStore.Call(
		uintptr(store),
		uintptr(windows.X509_ASN_ENCODING),
		uintptr(unsafe.Pointer(&certDer.Bytes[0])),
		uintptr(len(certDer.Bytes)),
		uintptr(4),
		0,
	)

	if r == 0 {
		return fmt.Errorf("调用CertAddEncodedCertificateToStore出错: %v", err)
	}
	return nil
}

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
