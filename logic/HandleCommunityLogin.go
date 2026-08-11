package logic

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"gf2gacha/config"
	"gf2gacha/logger"
	"gf2gacha/request"
	"gf2gacha/util"
	"strings"

	"github.com/pkg/errors"
)

const charSet = "01234567890123456789012345678901" +
	" !\"#$%&'()*+,-./" +
	"0123456789:;<=>?" +
	"@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_" +
	"`abcdefghijklmnopqrstuvwxyz{|}~"

var charToIdx [128]byte

func init() {
	for i := 0; i < len(charSet); i++ {
		charToIdx[charSet[i]] = byte(i)
	}
}

func CustomHash(password string) string {
	buf := make([]byte, len(password))
	for i := 0; i < len(password); i++ {
		buf[i] = charToIdx[password[i]]
	}

	hash := md5.Sum(buf)

	var sb strings.Builder
	sb.Grow(32)
	for i := 0; i < 4; i++ {
		wordBe := uint32(hash[i*4])<<24 |
			uint32(hash[i*4+1])<<16 |
			uint32(hash[i*4+2])<<8 |
			uint32(hash[i*4+3])
		sb.WriteString(fmt.Sprintf("%08x", wordBe))
	}
	return sb.String()
}

var key = []byte("a86a86^oH$04r6A1")

func LoginAESEncrypt(plaintext string) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	plainBytes := []byte(plaintext)

	// PKCS7 填充
	padLen := aes.BlockSize - len(plainBytes)%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	plainBytes = append(plainBytes, padding...)

	ciphertext := make([]byte, len(plainBytes))
	mode := cipher.NewCBCEncrypter(block, key)
	mode.CryptBlocks(ciphertext, plainBytes)

	// 转为 URL 安全的 base64
	b64 := base64.StdEncoding.EncodeToString(ciphertext)
	b64 = strings.ReplaceAll(b64, "+", "-")
	b64 = strings.ReplaceAll(b64, "/", "_")
	b64 = strings.TrimRight(b64, "=")
	return b64
}

func HandleCommunityLoginWithAccount(accountName string, passwd string) (webToken string, err error) {
	source := "phone"
	if strings.Contains(accountName, "@") {
		source = "mail"
	}

	encryptedAccount := LoginAESEncrypt(accountName)
	encryptedPasswd := LoginAESEncrypt(CustomHash(passwd))

	webToken, err = request.CommunityLoginWithAccount(encryptedAccount, encryptedPasswd, source)
	if err != nil {
		var respData request.CommonResponse
		if errors.As(err, &respData) && respData.Code == -1 {
			logger.Logger.Error(respData.Message)
			return "", errors.New(respData.Message)
		}
		return "", errors.WithStack(err)
	}
	return webToken, nil
}

func HandleCommunityLogin() (messageList []string, err error) {
	logInfo, err := util.GetLogInfo()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	AccountName := config.GetAccountName(logInfo.Uid)
	if AccountName == "" {
		return nil, errors.New("未配置账号")
	}

	Passwd := config.GetPasswd(logInfo.Uid)
	if Passwd == "" {
		return nil, errors.New("未配置密码")
	}

	webToken, err := HandleCommunityLoginWithAccount(AccountName, Passwd)
	if err != nil {
		return nil, err
	}
	logger.Logger.Infof("(UID: %s)登陆成功", logInfo.Uid)
	err = config.SetWebToken(logInfo.Uid, webToken)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return messageList, nil
}
