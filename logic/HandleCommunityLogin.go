package logic

import (
	"gf2gacha/config"
	"gf2gacha/encrypt"
	"gf2gacha/logger"
	"gf2gacha/request"
	"gf2gacha/util"
	"strings"

	"github.com/pkg/errors"
)

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

	source := "phone"
	if strings.Contains(AccountName, "@") {
		source = "mail"
	}
	// 加密账号和密码
	encryptedAccount := encrypt.LoginAESEncrypt(AccountName)
	encryptedPasswd := encrypt.LoginAESEncrypt(encrypt.CustomHash(Passwd))

	logger.Logger.Infof(encryptedAccount, encryptedPasswd)

	webToken, err := request.CommunityLoginWithAccount(encryptedAccount, encryptedPasswd, source)
	if err != nil {
		var respData request.CommonResponse
		if errors.As(err, &respData) {
			if respData.Code == -1 {
				logger.Logger.Errorf(respData.Message)
				return nil, errors.New(respData.Message)
			} else {
				return nil, errors.WithStack(err)
			}
		} else {
			return nil, errors.WithStack(err)
		}
	}
	logger.Logger.Infof("(UID: %s)登陆成功", logInfo.Uid)
	err = config.SetWebToken(logInfo.Uid, webToken)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return messageList, nil
}
