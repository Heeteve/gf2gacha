package logic

import (
	"gf2gacha/logger"
	"gf2gacha/model"
	"gf2gacha/preload"
	"strings"

	"github.com/pkg/errors"
)

func GetPoolInfo(uid string, poolType int64) (model.Pool, error) {
	localRecordList, err := GetLocalRecord(uid, poolType, 0)
	if err != nil {
		return model.Pool{}, errors.WithStack(err)
	}

	pool := model.Pool{PoolType: poolType}
	var isPreviousLose bool
	for _, storedRecord := range localRecordList {
		pool.GachaCount++
		pool.StoredCount++
		item := preload.ItemMap[storedRecord.ItemId]
		itemRank := item.Rank
		if poolType == 8 {
			itemRank = preload.ItemRankMap[storedRecord.PoolId][storedRecord.ItemId]
		}

		// 皮肤池只显示衣装记录
		if poolType == 9 {
			name := preload.LangMap[item.Name.Id]
			if itemRank == 5 {
				// 如果Name以“衣装·”开头
				if strings.HasPrefix(name, "衣装·") {
					pool.RecordList = append(pool.RecordList, model.DisplayRecord{
						Id:    item.Id,
						Name:  name,
						Lose:  false,
						Count: pool.StoredCount,
					})
				}
				pool.StoredCount = 0
				pool.Rank5Count++
			} else if itemRank == 4 {
				if strings.HasPrefix(name, "衣装·") {
					pool.RecordList = append(pool.RecordList, model.DisplayRecord{
						Id:    item.Id,
						Name:  name,
						Lose:  false,
						Count: pool.StoredCount,
					})
				}
				pool.Rank4Count++
			} else if itemRank == 3 {
				pool.Rank3Count++
			} else {
				logger.Logger.Warnf("未知的物品Rank poolType:%d poolId:%d itemId:%d", poolType, storedRecord.PoolId, storedRecord.ItemId)
			}
			continue
		}

		if itemRank == 5 {
			if isPreviousLose {
				pool.GuaranteesCount++
			}
			//检测是否歪
			var lose bool
			if upItemId, hasUp := preload.UpItemMap[storedRecord.PoolId]; hasUp && upItemId != storedRecord.ItemId {
				pool.LoseCount++
				lose = true
				isPreviousLose = true
			} else {
				isPreviousLose = false
			}

			pool.RecordList = append(pool.RecordList, model.DisplayRecord{
				Id:    item.Id,
				Name:  preload.LangMap[item.Name.Id],
				Lose:  lose,
				Count: pool.StoredCount,
			})

			pool.StoredCount = 0
			pool.Rank5Count++
		} else if itemRank == 4 {
			pool.Rank4Count++
		} else if itemRank == 3 {
			pool.Rank3Count++
		} else {
			logger.Logger.Warnf("未知的物品Rank poolType:%d poolId:%d itemId:%d", poolType, storedRecord.PoolId, storedRecord.ItemId)
		}
	}

	return pool, nil
}
