package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/utils"
)

func tokensImgTask(chainId int, tokenAddress string) {
	//defer time.AfterFunc(time.Duration(60*60*24)*time.Second, tokensImgTask)

	var count int
	sqlSel := oo.NewSqler().Table(consts.TbNameTokensImg).
		Where("chain_id", chainId).Where("token_address", tokenAddress).Count()
	err := oo.SqlGet(sqlSel, &count)
	if err != nil {
		oo.LogW("SQL failed. err: %v\n", err)
		return
	}

	if count == 0 {
		var platforms string
		switch chainId {
		case 1:
			platforms = "ethereum"
			break
		case 80001:
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
		for index := range resId {
			if resId[index].Platforms[platforms] == tokenAddress {
				imgStr := fmt.Sprintf(`https://api.coingecko.com/api/v3/coins/%s?localization=false&tickers=false&market_data=false&community_data=false&developer_data=false&sparkline=false`, resId[index].Id)
				resImg, err := utils.GetTokenImg(imgStr)
				if err != nil {
					oo.LogW("GetTokenImg failed error: %v", err)
					return
				}

				var m = make([]map[string]interface{}, 0)
				var v = make(map[string]interface{})
				v["chain_id"] = chainId
				v["token_address"] = tokenAddress
				v["thumb"] = resImg.Image.Thumb
				v["small"] = resImg.Image.Small
				v["large"] = resImg.Image.Large
				m = append(m, v)
				sqlIns := oo.NewSqler().Table(consts.TbNameTokensImg).Insert(m)
				err = oo.SqlExec(sqlIns)
				if err != nil {
					oo.LogW("SQL failed. err: %v\n", err)
					return
				}
				break
			}
		}
	}

}
