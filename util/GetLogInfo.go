package util

import (
	"fmt"
	"gf2gacha/encrypt"
	"gf2gacha/model"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"regexp"
)

func GetLogInfo() (logInfo model.LogInfo, err error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return model.LogInfo{}, errors.WithStack(err)
	}

	logPath := filepath.Join(userHome, "/AppData/LocalLow/SunBorn/少女前线2：追放/Player.log")
	logData, err := os.ReadFile(logPath)
	if err != nil {
		return model.LogInfo{}, errors.WithStack(err)
	}

	// logData补充读取./capture.log
	captureLogPath := filepath.Join("./capture.log")
	captureLogEncrypt, err := os.ReadFile(captureLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return model.LogInfo{}, errors.New("未找到捕获记录capture.log，官服请使用「捕获信息」按钮抓包")
		} else {
			return model.LogInfo{}, errors.WithStack(err)
		}
	}
	// 进行解密
	aesKey := encrypt.AesKey
	captureLogData, err := encrypt.Decrypt(captureLogEncrypt, aesKey)
	if err != nil {
		fmt.Printf("解密文件时出错: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("captureLogData: %v\n", captureLogData)
	logData = append(logData, captureLogData...)

	regexpGamePath, err := regexp.Compile(`\[Subsystems] Discovering subsystems at path (.+)/UnitySubsystems`)
	if err != nil {
		return model.LogInfo{}, errors.WithStack(err)
	}
	resultGamePath := regexpGamePath.FindSubmatch(logData)
	if len(resultGamePath) == 2 {
		logInfo.TablePath = filepath.Join(string(resultGamePath[1]), "LocalCache/Data/Table")
	} else {
		return model.LogInfo{}, errors.New("未在日志中找到游戏路径")
	}

	regexpUserInfo, err := regexp.Compile(`"access_token":"(.+?)".+"uid":(\d+)`)
	if err != nil {
		return model.LogInfo{}, errors.WithStack(err)
	}
	resultUserInfoList := regexpUserInfo.FindAllSubmatch(logData, -1)
	if len(resultUserInfoList) == 0 {
		return model.LogInfo{}, errors.New("未在日志中找到AccessToken或Uid,可能是最近一次游戏启动时未登录,官服请使用「捕获信息」按钮抓包")
	}
	resultUserInfo := resultUserInfoList[len(resultUserInfoList)-1]
	if len(resultUserInfo) == 3 {
		logInfo.AccessToken = string(resultUserInfo[1])
		logInfo.Uid = string(resultUserInfo[2])
	} else {
		return model.LogInfo{}, errors.New("未在日志中找到AccessToken或Uid,可能是最近一次游戏启动时未登录,官服请使用「捕获信息」按钮抓包")
	}

	regexpGachaUrl, err := regexp.Compile(`"gacha_record_url":"(.*?)"`)
	if err != nil {
		return model.LogInfo{}, errors.WithStack(err)
	}
	resultGachaUrlList := regexpGachaUrl.FindAllSubmatch(logData, -1)
	if len(resultGachaUrlList) == 0 {
		return model.LogInfo{}, errors.New("未在日志中找到抽卡链接,官服请使用「捕获信息」按钮抓包")
	}
	resultGachaUrl := resultGachaUrlList[len(resultGachaUrlList)-1]
	if len(resultGachaUrl) == 2 {
		logInfo.GachaUrl = string(resultGachaUrl[1])
	} else {
		return model.LogInfo{}, errors.New("未在日志中找到抽卡链接")
	}

	return logInfo, nil
}
