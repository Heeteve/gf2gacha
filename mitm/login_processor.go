package mitm

import (
	"fmt"
	"gf2gacha/config"
	"gf2gacha/util"
	"regexp"
)

// LoginProcessor 处理登录响应的处理器
type LoginProcessor struct {
	targetPattern *regexp.Regexp
}

// NewLoginProcessor 创建新的登录处理器
func NewLoginProcessor() *LoginProcessor {
	// 匹配 https://gf2-zoneinfo.sunborngame.com/<任意字符>login<任意字符>
	pattern := regexp.MustCompile(`https://gf2-zoneinfo\.sunborngame\.com/.*login.*`)
	return &LoginProcessor{
		targetPattern: pattern,
	}
}

// GetName 获取处理器名称
func (p *LoginProcessor) GetName() string {
	return "LoginProcessor"
}

// Match 检查URL是否与目标URL匹配
func (p *LoginProcessor) Match(url string) bool {
	return p.targetPattern.MatchString(url)
}

// Process 处理捕获到的数据
func (p *LoginProcessor) Process(capture *MitmCapture) (interface{}, error) {
	var resp LoginResponse
	if err := capture.ParseJSON(&resp); err != nil {
		return nil, fmt.Errorf("解析登录响应失败: %v", err)
	}

	if resp.Code != 0 || resp.Data.AccessToken == "" {
		return nil, fmt.Errorf("登录响应无效: %s", resp.Msg)
	}

	// 创建登录信息
	loginInfo := &LoginInfo{
		AccessToken: resp.Data.AccessToken,
		Uid:         fmt.Sprint(resp.Data.Uid),
		RawCapture:  capture,
	}

	// 保存配置和日志
	if err := p.saveLoginData(loginInfo); err != nil {
		fmt.Printf("保存登录数据时出现警告: %v\n", err)
	}

	return loginInfo, nil
}

// saveLoginData 保存登录数据到配置和日志
func (p *LoginProcessor) saveLoginData(loginInfo *LoginInfo) error {
	// 保存accessToken到config
	if err := config.SetAccessToken(loginInfo.Uid, loginInfo.AccessToken); err != nil {
		return fmt.Errorf("保存accessToken到config失败: %v", err)
	}

	// 删除旧的capture.log
	if err := util.DeleteCaptureLog(); err != nil {
		return fmt.Errorf("删除capture.log失败: %v", err)
	}

	// 将登录信息写入Player.log
	if err := util.WriteLoginInfoToLog(loginInfo.AccessToken, loginInfo.Uid); err != nil {
		return fmt.Errorf("写入capture.log失败: %v", err)
	}

	// TODO: 登录URL的抓取（因为抓取还要打开抽卡记录页面，懒得每次都抓了，等URL变了再说XD）
	util.WriteGachaUrlToLog("https://gf2-gacha-record.sunborngame.com/list")

	return nil
}
