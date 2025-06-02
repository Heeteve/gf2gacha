package logic

import (
	"fmt"
	"gf2gacha/model"
	"gf2gacha/request"
	"github.com/pkg/errors"
)

func GetCommunityExchangeList() ([]model.CommunityExchangeList, error) {
	exchangeListData, err := request.CommunityExchangeList("")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var exchangeList []model.CommunityExchangeList
	for _, exchangeItem := range exchangeListData.List {
		exchangeList = append(exchangeList, model.CommunityExchangeList{
			Id:   exchangeItem.ExchangeId,
			Name: fmt.Sprintf("『%s*%d』x%d", exchangeItem.ItemName, exchangeItem.ItemCount, exchangeItem.MaxExchangeCount),
		})
	}
	return exchangeList, nil
}
