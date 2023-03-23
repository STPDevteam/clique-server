package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
	"time"
)

func (svc *Service) doPush() {
	var err error
	var api *utils.PushAPI
	api, err = utils.NewPushAPI(
		svc.pushConfig.Endpoint,
		svc.pushConfig.Source,
		svc.pushConfig.ChannelAddress,
		svc.pushConfig.ChainId,
		svc.appConfig.SignMessagePriKey,
	)
	if err != nil {
		oo.LogD("doPush start failed: %v", err)
		return
	}

	var lastId = svc.pushConfig.StartId
	var page, pageSize uint64 = 1, 10
	for {
		var ret *utils.PageFeeds
		if ret, err = api.GetFeeds(page, pageSize); err == nil {
			for _, feed := range ret.Feeds {
				if feed.Payload.Recipients == api.GetChannel() {
					if vproof := strings.Split(feed.Payload.VerificationProof, "::uid::"); len(vproof) == 2 {
						if uids := strings.Split(vproof[1], "--"); len(uids) == 2 {
							lastId = uids[0]
							break
						}
					}
				}
			}
		}
		if lastId != "0" || ret.Itemcount < page*pageSize {
			break
		}
		page++
		// PUSH API rate limit?
		time.Sleep(time.Second)
	}

	oo.LogD("start PUSH ...")
	for {
		var data []models.EventHistoricalModel
		sqlStr := oo.NewSqler().Table(consts.TbNameEventHistorical).Where("id", ">", lastId).Limit(10).Select()
		if err = oo.SqlSelect(sqlStr, &data); err == nil && len(data) > 0 {
			for _, datum := range data {
				uid := fmt.Sprintf(
					"%d--%s",
					datum.Id,
					utils.Keccak256(fmt.Sprintf("%s-%s-%s", datum.BlockNumber, datum.TransactionHash, datum.LogIndex)),
				)
				payload := map[string]interface{}{
					"notification": map[string]string{
						"title": fmt.Sprintf("[MyClique] %s", datum.EventType),
						"body":  fmt.Sprintf("%s %s at tx: %s", datum.MessageSender, datum.EventType, datum.TransactionHash),
					},
					"data": map[string]string{
						"acta": "",
						"aimg": "",
						"amsg": "",
						"asub": "",
						"type": "1",
					},
					"recipients": api.GetChannelCAIPAddress(),
				}

				if err = api.SendNotification(uid, payload); err == nil {
					lastId = fmt.Sprintf("%d", datum.Id)
				}
			}
		}

		if err != nil {
			oo.LogW("error occurred err: %v", err)
		}
		time.Sleep(time.Second * 5)
	}
}
