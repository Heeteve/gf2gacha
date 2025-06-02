package util

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sys/windows"
	"math/big"
	"os"
	"syscall"
	"time"
	"unsafe"
)

func GenCA() error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "GF2Gacha Custom CA",
			Organization: []string{"GF2Gacha Custom Org"},
			Country:      []string{"CN"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(20, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return err
	}
	certOut, err := os.Create("ca.crt")
	if err != nil {
		return err
	}
	defer certOut.Close()
	err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return err
	}

	keyOut, err := os.Create("ca.key")
	if err != nil {
		return err
	}
	defer keyOut.Close()
	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	if err != nil {
		return err
	}

	return nil
}

var (
	crypt32                              = syscall.NewLazyDLL("crypt32.dll")
	procCertOpenStore                    = crypt32.NewProc("CertOpenStore")
	procCertCloseStore                   = crypt32.NewProc("CertCloseStore")
	procCertAddEncodedCertificateToStore = crypt32.NewProc("CertAddEncodedCertificateToStore")
)

func IsTrustedCA() bool {
	certPem, err := os.ReadFile("ca.crt")
	if err != nil {
		fmt.Println(err)
		return false
	}
	certDer, _ := pem.Decode(certPem)
	if certDer == nil {
		fmt.Println("pem解码失败")
		return false
	}
	cert, err := x509.ParseCertificate(certDer.Bytes)
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, err = cert.Verify(x509.VerifyOptions{})
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

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

func InstallCA() error {
	certPem, err := os.ReadFile("ca.crt")
	if err != nil {
		return err
	}
	certDer, _ := pem.Decode(certPem)
	if certDer == nil {
		return errors.New("pem解码失败")
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

func ParseCA() (*tls.Certificate, error) {
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		return nil, err
	}

	caKey, err := os.ReadFile("ca.key")
	if err != nil {
		return nil, err
	}

	parsedCert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return nil, err
	}
	if parsedCert.Leaf, err = x509.ParseCertificate(parsedCert.Certificate[0]); err != nil {
		return nil, err
	}
	return &parsedCert, nil
}
