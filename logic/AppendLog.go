package logic

import (
	_ "embed"
	"fmt"
	"gf2gacha/encrypt"
	"os"
	"path/filepath"
)

func AppendLog(accessToken string, uid int64, gachaUrl string) error {
	logPath := filepath.Join("./capture.log")
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 加密数据
	aesKey := encrypt.AesKey
	plainData := fmt.Sprintf(`{"access_token":"%s","uid":%d,"gacha_record_url":"%s"}`, accessToken, uid, gachaUrl)
	encryptedData, err := encrypt.Encrypt([]byte(plainData), aesKey)
	if err != nil {
		return err
	}

	// 写入加密数据
	err = os.WriteFile(logPath, encryptedData, 0644);
	if err != nil {
		return err
	}
	return nil
}
