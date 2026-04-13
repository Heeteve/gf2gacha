package request

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type CommunityLoginData struct {
	Account struct {
		Token      string `json:"token"`
		Uid        int64  `json:"uid"`
		PlatformId int64  `json:"platform_id"`
		ChannelId  int64  `json:"channel_id"`
	} `json:"account"`
}

func CommunityLogin(gameToken string) (webToken string, err error) {
	apiUrl := `https://gf2-bbs-api.sunborngame.com/login/game_skip`
	params := map[string]interface{}{
		"game_token": gameToken,
	}

	dataBytes, err := CommunityPost(apiUrl, params, "")
	if err != nil {
		return "", errors.WithStack(err)
	}

	var data CommunityLoginData
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return data.Account.Token, nil
}

func CommunityLoginWithAccount(accountName string, password string) (webToken string, err error) {
	apiUrl := `https://gf2-bbs-api.sunborngame.com/login/account`
	source := "phone"
	if strings.Contains(accountName, "@") {
		source = "mail"
	}
	params := map[string]interface{}{
		"account_name": accountName,
		"passwd":       password,
		"source":       source,
	}

	dataBytes, err := CommunityPost(apiUrl, params, "")
	if err != nil {
		return "", errors.WithStack(err)
	}

	var data CommunityLoginData
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return data.Account.Token, nil
}
