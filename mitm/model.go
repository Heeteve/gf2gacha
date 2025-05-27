package mitm

import (
	"encoding/json"
	"fmt"
)

// MitmCapture 单次抓取的结果
type MitmCapture struct {
	URL      string            // 请求的 URL
	Headers  map[string]string // 响应头
	Response []byte            // 原始响应
	JsonData interface{}       // 解析后的 JSON 数据，类型由具体业务决定
}

// MitmProcessor 定义抓包处理器接口
type MitmProcessor interface {
	// Match 检查是否需要处理该 URL
	Match(url string) bool
	// Process 处理抓包数据，并返回处理结果
	Process(capture *MitmCapture) (interface{}, error)
	// GetName 获取处理器名称，用于标识
	GetName() string
}

// LoginResponse 登录接口返回的数据结构
type LoginResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		AccessToken string `json:"access_token"`
		Uid         int    `json:"uid"`
	} `json:"data"`
}

// LoginInfo 提取的登录信息
type LoginInfo struct {
	AccessToken string
	Uid         string
	RawCapture  *MitmCapture
}

// ParseJSON 通用JSON解析方法
func (c *MitmCapture) ParseJSON(target interface{}) error {
	if err := json.Unmarshal(c.Response, target); err != nil {
		return fmt.Errorf("JSON解析失败: %v", err)
	}
	c.JsonData = target
	return nil
}
