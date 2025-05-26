package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeToPlayerLog(content string) error {

	logPath := filepath.Join("./capture.log")

	// 确保目录存在
	logDir := filepath.Dir(logPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 以追加模式打开文件
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开capture.log失败: %v", err)
	}
	defer file.Close()

	if _, err := file.Write([]byte(content + "\n")); err != nil {
		return fmt.Errorf("写入capture.log失败: %v", err)
	}

	return nil
}

func WriteLoginInfoToLog(accessToken string, uid string) error {
	// 创建要写入的登录信息
	loginInfo := fmt.Sprintf(`{"access_token":"%s","uid":%s}`, accessToken, uid)
	return writeToPlayerLog(loginInfo)
}

func WriteGachaUrlToLog(gachaUrl string) error {
	// 创建要写入的抽卡链接
	gachaInfo := fmt.Sprintf(`{"gacha_record_url":"%s"}`, gachaUrl)
	return writeToPlayerLog(gachaInfo)
}

func DeleteCaptureLog() error {
	logPath := filepath.Join("./capture.log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return nil // 文件不存在，直接返回
	}
	if err := os.Remove(logPath); err != nil {
		return fmt.Errorf("删除capture.log失败: %v", err)
	}
	return nil
}