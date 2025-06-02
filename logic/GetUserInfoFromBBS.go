package logic

import "gf2gacha/request"

func GetUserInfoFromBBS(accessToken string) (*request.CommunityUserInfoData, error) {
	webToken, err := request.CommunityLogin(accessToken)
	if err != nil {
		return nil, err
	}
	userInfo, err := request.CommunityUserInfo(webToken)
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}
