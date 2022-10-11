package controller

import (
	"time"
)

func (svc *Service) getV1Proposal() {
	defer time.AfterFunc(time.Duration(1)*time.Second, svc.getV1Proposal)

	//res, err := utils.GetV1LastBlockNumber(svc.appConfig.ApiV1BlockUrl)
	//if err != nil {
	//	oo.LogW("GetV1LastBlockNumber failed error: %v", err)
	//	return
	//}
}
