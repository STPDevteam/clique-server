package handlers

import (
	"encoding/json"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
)

func checkLogin(sign *models.SignData) (ret bool) {
	ret, errSign := utils.CheckPersonalSign(consts.SignMessagePrefix, sign.Account, sign.Signature)
	if errSign != nil {
		oo.LogD("signMessage err %v, signature: %v", errSign, sign.Signature)
		return false
	}
	if !ret {
		oo.LogD("check Sign fail")
		return false
	}
	return true
}

func checkAirdropAdminAndTimestamp(sign *models.AirdropAdminSignData) (ret bool) {
	ret, errSign := utils.CheckPersonalSign(sign.Message, sign.Account, sign.Signature)
	if errSign != nil {
		oo.LogD("signMessage err %v", errSign)
		return false
	}

	var data models.AdminMessage
	err := json.Unmarshal([]byte(sign.Message), &data)
	if err != nil {
		oo.LogD("signMessage Unmarshal err %v", errSign)
		return false
	}

	if !utils.CheckAdminSignMessageTimestamp(data.Expired) {
		oo.LogD("signMessage CheckAdminSignMessageTimestamp err %v", errSign)
		return false
	}

	//if data.Type == "airdrop2" {
	//	root, err := merkelTreeRoot(sign.Array)
	//	log.Println(fmt.Sprintf("rootStr: %s", root))
	//	if err != nil || root != data.Root {
	//		oo.LogD("signMessage err rootMe:%v.root:%v", root, data.Root)
	//		return false
	//	}
	//}

	var count int
	var sqlSql string
	if data.Type == "airdrop1" {
		sqlSql = oo.NewSqler().Table(consts.TbNameAdmin).Where("chain_id", sign.ChainId).Where("dao_address", sign.DaoAddress).
			Where("account", sign.Account).Where("account_level='superAdmin' OR account_level='admin'").Count()
	} else if data.Type == "airdrop2" || data.Type == "airdropDownload" {
		sqlSql = oo.NewSqler().Table(consts.TbNameAirdrop).Where("id", sign.AirdropId).Where("creator", sign.Account).Count()
	}
	err = oo.SqlGet(sqlSql, &count)
	if err != nil || count == 0 {
		return false
	}

	if !ret {
		oo.LogD("check Sign fail")
		return false
	}
	return true
}

func checkAccountJoinOrQuit(data *models.JoinDaoWithSignParam) (ret bool) {
	message := fmt.Sprintf(`%d,%s,%s,%d`, data.Params.ChainId, data.Params.DaoAddress, data.Params.JoinSwitch, data.Params.Timestamp)
	ret, errSign := utils.CheckPersonalSign(message, data.Sign.Account, data.Sign.Signature)
	if errSign != nil {
		oo.LogD("signMessage err %v", errSign)
		return false
	}

	if !utils.CheckAdminSignMessageTimestamp(data.Params.Timestamp) {
		oo.LogD("signMessage CheckAdminSignMessageTimestamp err %v", errSign)
		return false
	}

	if !ret {
		oo.LogD("check Sign fail")
		return false
	}
	return true
}

//func checkAdminOrMember(data models.SignDataForTask) (role string, ret bool) {
//	message := fmt.Sprintf(`%d,%s,%s,%d`, data.ChainId, data.DaoAddress, data.Account, data.Timestamp)
//	ret, err := utils.CheckPersonalSign(message, data.Account, data.Signature)
//	if err != nil {
//		oo.LogW("signMessage err:%v", err)
//		return "", false
//	}
//
//	if !ret {
//		oo.LogW("check Sign failed.")
//		return "", false
//	}
//
//	if !utils.CheckAdminSignMessageTimestamp(data.Timestamp) {
//		oo.LogW("signMessage deadline.")
//		return "", false
//	}
//
//	jobs, err := db.GetTbJobs(
//		o.W("chain_id", data.ChainId),
//		o.W("dao_address", data.DaoAddress),
//		o.W("account", data.Account))
//	if err != nil {
//		oo.LogW("SQL err:%v", err)
//		return "", false
//	}
//
//	return jobs.Job, true
//}

func IsSuperAdmin(chainId int64, daoAddress, account string) (b bool) {
	jobData, err := db.GetTbJobs(
		o.W("chain_id", chainId),
		o.W("dao_address", daoAddress),
		o.W("account", account))
	if err != nil {
		oo.LogW("SQL err:%v", err)
		return false
	}
	if jobData.Job != consts.Jobs_A_superAdmin {
		return false
	}

	return true
}

func IsAboveAdmin(chainId int64, daoAddress, account string) (role string, b bool) {
	jobData, err := db.GetTbJobs(
		o.W("chain_id", chainId),
		o.W("dao_address", daoAddress),
		o.W("account", account))
	if err != nil {
		oo.LogW("SQL err:%v", err)
		return "", false
	}
	if jobData.Job != consts.Jobs_A_superAdmin && jobData.Job != consts.Jobs_B_admin {
		return "", false
	}

	return jobData.Job, true
}

func IsTaskAssign(taskId int64, account string) (b bool) {
	task, err := db.GetTbTask(o.W("id", taskId))
	if err != nil || task.AssignAccount == "" || strings.ToLower(task.AssignAccount) != strings.ToLower(account) {
		return false
	}

	return true
}
