package handlers

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/db"
	"stp_dao_v2/db/o"
	"stp_dao_v2/models"
	"time"
)

func PageTbTeamSpaces(order string, page ReqPagination, w ...[][]interface{}) (list []models.ResTeamSpacesList, total int64, err error) {
	var data []db.TbTeamSpaces
	sqler := o.Pre(consts.TbTeamSpaces, w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Order(order).Limit(page.Limit).Offset(page.Offset).Select(), &data)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	list = make([]models.ResTeamSpacesList, 0)
	for i := range data {
		ls := data[i]

		var avatarCreator, nicknameCreator string
		if ls.Creator != "" {
			account, err := db.GetTbAccountModel(o.W("account", ls.Creator))
			if err != nil && err != oo.ErrNoRows {
				return nil, 0, err
			}
			avatarCreator = account.AccountLogo.String
			nicknameCreator = account.Nickname.String
		}
		var avatarLastEditBy, nicknameLastEditBy string
		if ls.LastEditBy != "" {
			account, err := db.GetTbAccountModel(o.W("account", ls.LastEditBy))
			if err != nil && err != oo.ErrNoRows {
				return nil, 0, err
			}
			avatarLastEditBy = account.AccountLogo.String
			nicknameLastEditBy = account.Nickname.String
		}

		list = append(list, models.ResTeamSpacesList{
			TeamSpacesId:       ls.Id,
			ChainId:            ls.ChainId,
			DaoAddress:         ls.DaoAddress,
			Creator:            ls.Creator,
			AvatarCreator:      avatarCreator,
			NicknameCreator:    nicknameCreator,
			Title:              ls.Title,
			Url:                ls.Url,
			LastEditTime:       ls.LastEditTime,
			LastEditBy:         ls.LastEditBy,
			AvatarLastEditBy:   avatarLastEditBy,
			NicknameLastEditBy: nicknameLastEditBy,
			Access:             ls.Access,
		})
	}
	return list, total, nil
}

func PageTbTask(order string, page ReqPagination, w ...[][]interface{}) (list []models.ResTaskList, total int64, err error) {
	var data []db.TbTask
	sqler := o.Pre(consts.TbTask, w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Order(order).Limit(page.Limit).Offset(page.Offset).Select(), &data)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	list = make([]models.ResTaskList, 0)
	for i := range data {
		ls := data[i]

		var avatar, nickname string
		if ls.AssignAccount != "" {
			account, err := db.GetTbAccountModel(o.W("account", ls.AssignAccount))
			if err != nil && err != oo.ErrNoRows {
				return nil, 0, err
			}
			avatar = account.AccountLogo.String
			nickname = account.Nickname.String
		}

		list = append(list, models.ResTaskList{
			SpacesId:       ls.SpacesId,
			TaskId:         ls.Id,
			TaskName:       ls.TaskName,
			Deadline:       ls.Deadline,
			Priority:       ls.Priority,
			AssignAccount:  ls.AssignAccount,
			AssignAvatar:   avatar,
			AssignNickname: nickname,
			ProposalId:     ls.ProposalId,
			Reward:         ls.Reward,
			Status:         ls.Status,
			Weight:         ls.Weight,
		})
	}
	return list, total, nil
}

func PageTbJobs(order string, page ReqPagination, w ...[][]interface{}) (list []models.ResJobsList, total int64, err error) {
	var data []db.TbJobs
	sqler := o.Pre(consts.TbJobs, w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Order(order).Limit(page.Limit).Offset(page.Offset).Select(), &data)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	list = make([]models.ResJobsList, 0)
	for i := range data {
		ls := data[i]

		var avatar, nickname, twitter, discord, youtube, opensea string
		if ls.Account != "" {
			account, err := db.GetTbAccountModel(o.W("account", ls.Account))
			if err != nil && err != oo.ErrNoRows {
				return nil, 0, err
			}
			avatar = account.AccountLogo.String
			nickname = account.Nickname.String
			twitter = account.Twitter.String
			discord = account.Discord.String
			youtube = account.Youtube.String
			opensea = account.Opensea.String
		}

		list = append(list, models.ResJobsList{
			JobId:      ls.Id,
			ChainId:    ls.ChainId,
			DaoAddress: ls.DaoAddress,
			Account:    ls.Account,
			Jobs:       ls.Job,
			Avatar:     avatar,
			Nickname:   nickname,
			Twitter:    twitter,
			Discord:    discord,
			Youtube:    youtube,
			Opensea:    opensea,
		})
	}
	return list, total, nil
}

func PageTbJobsApply(order string, page ReqPagination, w ...[][]interface{}) (list []models.ResJobsApplyList, total int64, err error) {
	var data []db.TbJobsApply
	sqler := o.Pre(consts.TbJobsApply, w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Order(order).Limit(page.Limit).Offset(page.Offset).Select(), &data)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	list = make([]models.ResJobsApplyList, 0)
	for i := range data {
		ls := data[i]

		var avatar, nickname string
		if ls.Account != "" {
			account, err := db.GetTbAccountModel(o.W("account", ls.Account))
			if err != nil && err != oo.ErrNoRows {
				return nil, 0, err
			}
			avatar = account.AccountLogo.String
			nickname = account.Nickname.String
		}

		createAt, _ := time.Parse("2006-01-02 15:04:05", ls.CreateTime)
		list = append(list, models.ResJobsApplyList{
			ApplyId:    ls.Id,
			ChainId:    ls.ChainId,
			DaoAddress: ls.DaoAddress,
			Account:    ls.Account,
			Avatar:     avatar,
			Nickname:   nickname,
			ApplyRole:  ls.ApplyRole,
			ApplyTime:  createAt.Unix(),
			Message:    ls.Message,
		})
	}
	return list, total, nil
}

func PageTbAccountTopList(order string, page ReqPagination, w ...[][]interface{}) (list []models.ResAccountTopList, total int64, err error) {
	var data []db.TbAccountModel
	sqler := o.Pre(consts.TbNameAccount, w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Order(order).Limit(page.Limit).Offset(page.Offset).Select(), &data)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	list = make([]models.ResAccountTopList, 0)
	for i := range data {
		ls := data[i]

		list = append(list, models.ResAccountTopList{
			Account:  ls.Account,
			Avatar:   ls.AccountLogo.String,
			Nickname: ls.Nickname.String,
			FansNum:  ls.FansNum,
		})
	}
	return list, total, nil
}

func PageTbSBT(order string, page ReqPagination, w ...[][]interface{}) (list []models.ResSBTList, total int64, err error) {
	var data []db.TbSBT
	sqler := o.Pre(consts.TbSBT, w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Order(order).Limit(page.Limit).Offset(page.Offset).Select(), &data)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	list = make([]models.ResSBTList, 0)
	for i := range data {
		ls := data[i]

		dao, err := db.GetTbDao(o.W("chain_id", ls.ChainId), o.W("dao_address", ls.DaoAddress))
		if err != nil {
			return nil, 0, err
		}

		list = append(list, models.ResSBTList{
			SBTId:        ls.Id,
			ChainId:      ls.ChainId,
			DaoAddress:   ls.DaoAddress,
			DaoName:      dao.DaoName,
			DaoLogo:      dao.DaoLogo,
			TokenChainId: ls.TokenChainId,
			FileUrl:      ls.FileUrl,
			ItemName:     ls.ItemName,
			Symbol:       ls.Symbol,
			StartTime:    ls.StartTime,
			EndTime:      ls.EndTime,
			Status:       ls.Status,
		})
	}
	return list, total, nil
}

func PageTbSBTClaim(order string, page ReqPagination, w ...[][]interface{}) (list []models.ResSBTClaimList, total int64, err error) {
	sqler := o.Pre(fmt.Sprintf(`%s AS a`, consts.TbSBTClaim), w)
	sqlCopy := *sqler
	err = oo.SqlGet(sqlCopy.Count(), &total)
	if err == nil {
		sqlCopy = *sqler
		err = oo.SqlSelect(sqlCopy.Join(fmt.Sprintf(`%s AS b on a.account = b.account`, consts.TbNameAccount)).
			Order(order).Limit(page.Limit).Offset(page.Offset).
			Select("a.account, IFNULL(b.account_logo,'') AS account_logo"), &list)
	}
	if err != nil {
		oo.LogW("sqler:%s", sqler)
		return nil, 0, err
	}

	return list, total, nil
}
