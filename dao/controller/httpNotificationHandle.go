package controller

import (
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"net/http"
	"stp_dao_v2/consts"
	"stp_dao_v2/models"
	"strconv"
)

// @Summary notification list
// @Tags notification
// @version 0.0.1
// @description notification list
// @Produce json
// @Param account query string true "account address"
// @Param offset query  int true "offset,page"
// @Param count query  int true "count,page"
// @Success 200 {object} models.ResNotificationPage
// @Router /stpdao/v2/notification/list [get]
func httpNotificationList(c *gin.Context) {
	accountParam := c.Query("account")
	count := c.Query("count")
	offset := c.Query("offset")
	countParam, _ := strconv.Atoi(count)
	offsetParam, _ := strconv.Atoi(offset)

	var accountEntities []models.NotificationAccountModel
	sqler := oo.NewSqler().Table(consts.TbNameNotificationAccount).Where("account", accountParam)

	var total uint64
	sqlCopy := *sqler
	sqlStr := sqlCopy.Count()
	err := oo.SqlGet(sqlStr, &total)
	if err == nil {
		sqlCopy = *sqler
		sqlStr = sqlCopy.Order("notification_time DESC").Limit(countParam).Offset(offsetParam).Select()
		err = oo.SqlSelect(sqlStr, &accountEntities)
	}
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	var data = make([]models.ResNotification, 0)
	for index := range accountEntities {
		var notificationEntities []models.NotificationModel
		sqlSel := oo.NewSqler().Table(consts.TbNameNotification).Where("id", accountEntities[index].NotificationId).Select()
		err = oo.SqlSelect(sqlSel, &notificationEntities)
		if err != nil {
			oo.LogW("SQL err: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    500,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}

		dataIndex := notificationEntities[0]

		var info = make([]models.NotificationInfo, 0)
		if dataIndex.Types == consts.TypesNameNewProposal {
			info = append(info, models.NotificationInfo{
				ChainId:      dataIndex.ChainId,
				DaoAddress:   dataIndex.DaoAddress,
				DaoLogo:      dataIndex.DaoLogo,
				DaoName:      dataIndex.DaoName,
				ProposalId:   dataIndex.ActivityId,
				ProposalName: dataIndex.ActivityName,
			})
		} else if dataIndex.Types == consts.TypesNameAirdrop {
			info = append(info, models.NotificationInfo{
				ChainId:      dataIndex.ChainId,
				DaoAddress:   dataIndex.DaoAddress,
				DaoLogo:      dataIndex.DaoLogo,
				DaoName:      dataIndex.DaoName,
				ActivityId:   dataIndex.ActivityId,
				ActivityName: dataIndex.ActivityName,
			})
		} else if dataIndex.Types == consts.TypesNamePublicSale {

		}

		data = append(data, models.ResNotification{
			Account:          accountEntities[index].Account,
			AlreadyRead:      accountEntities[index].AlreadyRead,
			NotificationId:   accountEntities[index].NotificationId,
			NotificationTime: dataIndex.NotificationTime,
			Types:            dataIndex.Types,
			Info:             info[0],
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResNotificationPage{
			List:  data,
			Total: total,
		},
	})
}

// @Summary notification read
// @Tags notification
// @version 0.0.1
// @description notification read
// @Produce json
// @Param request body models.NotificationReadParam true "request"
// @Success 200 {object} models.Response
// @Router /stpdao/v2/notification/read [post]
func httpNotificationRead(c *gin.Context) {
	var params models.NotificationReadParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		oo.LogW("%v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters.",
		})
		return
	}

	var sqlUp string
	if params.ReadAll {
		sqlUp = fmt.Sprintf(`UPDATE %s SET already_read=%t WHERE account='%s'`,
			consts.TbNameNotificationAccount,
			true,
			params.Account,
		)
	} else {
		sqlUp = fmt.Sprintf(`UPDATE %s SET already_read=%t WHERE account='%s' AND notification_id=%d`,
			consts.TbNameNotificationAccount,
			true,
			params.Account,
			params.NotificationId,
		)
	}
	err = oo.SqlExec(sqlUp)
	if err != nil {
		oo.LogW("SQL err: %v", err)
		c.JSON(http.StatusInternalServerError, models.Response{
			Code:    500,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Code:    http.StatusOK,
		Message: "ok",
		Data: models.ResResult{
			Success: true,
		},
	})
}
