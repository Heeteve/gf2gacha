package model

import (
	"gf2gacha/logger"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
)

var Engine *xorm.Engine

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite", "./gf2gacha.db?_pragma=busy_timeout(2000)")
	if err != nil {
		logger.Logger.Panic(err)
	}
}
