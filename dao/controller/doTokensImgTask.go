package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"time"
)

func tokensImgTask() {
	defer time.AfterFunc(time.Duration(60*60*24)*time.Second, tokensImgTask)

	var entities []models.DaoModel
	sqlNeed := oo.NewSqler().Table(consts.TbNameDao).Select("token_chain_id,token_address")
	err := oo.SqlSelect(sqlNeed, &entities)
	if err != nil {
		oo.LogW("SQL failed. err: %v\n", err)
		return
	}

	for indexToken := range entities {
		tokenChainId := entities[indexToken].TokenChainId
		tokenAddress := entities[indexToken].TokenAddress

		var platforms string
		switch tokenChainId {
		case 1:
			platforms = "ethereum"
			break
		case 56:
			platforms = "binance-smart-chain"
			break
		case 137:
			platforms = "polygon-pos"
			break
		default:
			platforms = "Undefined"
			break
		}

		resId, err := utils.GetTokensId("https://api.coingecko.com/api/v3/coins/list?include_platform=true")
		fmt.Println(err)
		if err != nil {
			oo.LogW("GetTokensId failed error: %v", err)
			return
		}
		for indexId := range resId {
			if resId[indexId].Platforms[platforms] == tokenAddress {
				imgStr := fmt.Sprintf(`https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=false&community_data=false&developer_data=false&sparkline=false`, resId[indexId].Id)
				resImg, err := utils.GetTokenImg(imgStr)
				if err != nil {
					oo.LogW("GetTokenImg failed error: %v", err)
					return
				}

				sqlIns := fmt.Sprintf(`REPLACE INTO %s (token_chain_id,token_address,thumb,small,large) VALUES (%d,'%s','%s','%s','%s')`,
					consts.TbNameTokensImg,
					tokenChainId,
					tokenAddress,
					resImg.Image.Thumb,
					resImg.Image.Small,
					resImg.Image.Large,
				)
				err = oo.SqlExec(sqlIns)
				if err != nil {
					oo.LogW("SQL failed. err: %v\n", err)
					return
				}
				break
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}

}
