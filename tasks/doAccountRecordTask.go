package tasks

import (
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"time"
)

func UpdateAccountRecord() {
	defer time.AfterFunc(time.Duration(60)*time.Second, UpdateAccountRecord)

	var entities []db.TbAccountRecordModel
	sqlSel := oo.NewSqler().Table(consts.TbNameAccountRecord).Where("update_bool", 1).Select()
	err := oo.SqlSelect(sqlSel, &entities)
	if err != nil {
		oo.LogW("query SQL account record failed. err:%v", err)
		return
	}

	if len(entities) != 0 {
		for index := range entities {

			if entities[index].Types == consts.EvCreateERC20 || entities[index].Types == consts.EvClaimReserve {
				var avatar string
				sqlSel = oo.NewSqler().Table(consts.TbNameTokensImg).Where("token_chain_id", entities[index].ChainId).
					Where("token_address", entities[index].Address).Select("own_img")
				err = oo.SqlGet(sqlSel, &avatar)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
				if avatar == "" || len(avatar) == 0 {
					return
				}
				var v = make(map[string]interface{})
				v["avatar"] = avatar
				v["update_bool"] = 0
				var types string
				if entities[index].Types == consts.EvCreateERC20 {
					types = consts.EvCreateERC20
				} else if entities[index].Types == consts.EvClaimReserve {
					types = consts.EvClaimReserve
				}
				sqlUp := oo.NewSqler().Table(consts.TbNameAccountRecord).Where("types", types).
					Where("chain_id", entities[index].ChainId).Where("address", entities[index].Address).Update(v)
				err = oo.SqlExec(sqlUp)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
			}

			if entities[index].Types == consts.EvCreateDao || entities[index].Types == consts.EvAdmin || entities[index].Types == consts.EvOwnershipTransferred {
				var daoEntity []db.TbDaoModel
				sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].Address).Select()
				err = oo.SqlSelect(sqlSel, &daoEntity)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
				if len(daoEntity) == 0 {
					return
				}
				if daoEntity[0].DaoLogo == "" || daoEntity[0].DaoName == "" {
					return
				}
				var v = make(map[string]interface{})
				v["avatar"] = daoEntity[0].DaoLogo
				v["dao_name"] = daoEntity[0].DaoName
				v["update_bool"] = 0
				var types string
				if entities[index].Types == consts.EvCreateDao {
					types = consts.EvCreateDao
				} else if entities[index].Types == consts.EvAdmin {
					types = consts.EvAdmin
				} else if entities[index].Types == consts.EvOwnershipTransferred {
					types = consts.EvOwnershipTransferred
				}
				sqlUp := oo.NewSqler().Table(consts.TbNameAccountRecord).Where("types", types).
					Where("chain_id", entities[index].ChainId).Where("address", entities[index].Address).Update(v)
				err = oo.SqlExec(sqlUp)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
			}

			if entities[index].Types == consts.EvCreateProposal || entities[index].Types == consts.EvCancelProposal || entities[index].Types == consts.EvVote {
				var daoEntity []db.TbDaoModel
				sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].Address).Select()
				err = oo.SqlSelect(sqlSel, &daoEntity)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
				if len(daoEntity) == 0 {
					return
				}
				if daoEntity[0].DaoLogo == "" || daoEntity[0].DaoName == "" {
					return
				}
				var proposalTitle string
				sqlSel = oo.NewSqler().Table(consts.TbNameProposal).Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].Address).Where("proposal_id", entities[index].ActivityId).Select("title")
				err = oo.SqlGet(sqlSel, &proposalTitle)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
				var v = make(map[string]interface{})
				v["avatar"] = daoEntity[0].DaoLogo
				v["dao_name"] = daoEntity[0].DaoName
				v["titles"] = proposalTitle
				v["update_bool"] = 0
				var types string
				if entities[index].Types == consts.EvCreateProposal {
					types = consts.EvCreateProposal
				} else if entities[index].Types == consts.EvCancelProposal {
					types = consts.EvCancelProposal
				} else if entities[index].Types == consts.EvVote {
					types = consts.EvVote
				}
				sqlUp := oo.NewSqler().Table(consts.TbNameAccountRecord).Where("types", types).
					Where("chain_id", entities[index].ChainId).Where("address", entities[index].Address).
					Where("activity_id", entities[index].ActivityId).Update(v)
				err = oo.SqlExec(sqlUp)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
			}

			if entities[index].Types == consts.EvCreateAirdrop || entities[index].Types == consts.EvSettleAirdrop || entities[index].Types == consts.EvClaimed {
				var daoEntity []db.TbDaoModel
				sqlSel = oo.NewSqler().Table(consts.TbNameDao).Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].Address).Select()
				err = oo.SqlSelect(sqlSel, &daoEntity)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
				if len(daoEntity) == 0 {
					return
				}
				if daoEntity[0].DaoLogo == "" || daoEntity[0].DaoName == "" {
					return
				}
				var airdropTitle string
				sqlSel = oo.NewSqler().Table(consts.TbNameAirdrop).Where("chain_id", entities[index].ChainId).
					Where("dao_address", entities[index].Address).Where("id", entities[index].ActivityId).Select("title")
				err = oo.SqlGet(sqlSel, &airdropTitle)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
				var v = make(map[string]interface{})
				v["avatar"] = daoEntity[0].DaoLogo
				v["dao_name"] = daoEntity[0].DaoName
				v["titles"] = airdropTitle
				v["update_bool"] = 0
				var types string
				if entities[index].Types == consts.EvCreateAirdrop {
					types = consts.EvCreateAirdrop
				} else if entities[index].Types == consts.EvSettleAirdrop {
					types = consts.EvSettleAirdrop
				} else if entities[index].Types == consts.EvClaimed {
					types = consts.EvClaimed
				}
				sqlUp := oo.NewSqler().Table(consts.TbNameAccountRecord).Where("types", types).
					Where("chain_id", entities[index].ChainId).Where("address", entities[index].Address).
					Where("activity_id", entities[index].ActivityId).Update(v)
				err = oo.SqlExec(sqlUp)
				if err != nil {
					oo.LogW("SQL err:%v", err)
					return
				}
			}

		}
	}
}
