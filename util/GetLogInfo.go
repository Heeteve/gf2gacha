package util

import (
	"gf2gacha/model"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
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
	captureLogData, err := os.ReadFile(captureLogPath)
	if err != nil {
		if os.IsNotExist(err) {
			return model.LogInfo{}, errors.New("未找到capture.log文件。请先开始抓取，然后启动游戏点击登录")
		} else {
			return model.LogInfo{}, errors.WithStack(err)
		}
	}
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
		return model.LogInfo{}, errors.New("未在日志中找到AccessToken或Uid,可能是最近一次游戏启动时未登录")
	}
	resultUserInfo := resultUserInfoList[len(resultUserInfoList)-1]
	if len(resultUserInfo) == 3 {
		logInfo.AccessToken = string(resultUserInfo[1])
		logInfo.Uid = string(resultUserInfo[2])
	} else {
		return model.LogInfo{}, errors.New("未在日志中找到AccessToken或Uid,可能是最近一次游戏启动时未登录")
	}

	regexpGachaUrl, err := regexp.Compile(`"gacha_record_url":"(.*?)"`)
	if err != nil {
		return model.LogInfo{}, errors.WithStack(err)
	}
	resultGachaUrlList := regexpGachaUrl.FindAllSubmatch(logData, -1)
	if len(resultGachaUrlList) == 0 {
		return model.LogInfo{}, errors.New("未在日志中找到抽卡链接")
	}
	resultGachaUrl := resultGachaUrlList[len(resultGachaUrlList)-1]
	if len(resultGachaUrl) == 2 {
		logInfo.GachaUrl = string(resultGachaUrl[1])
	} else {
		return model.LogInfo{}, errors.New("未在日志中找到抽卡链接")
	}

	return logInfo, nil
}
