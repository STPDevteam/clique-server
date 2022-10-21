package controller

import (
	_ "crypto/ecdsa"
	"encoding/json"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/urfave/cli"
	"os"
	"path"
	"path/filepath"
	"stp_dao_v2/config"
	"stp_dao_v2/consts"
	_ "stp_dao_v2/dao/docs"
	"stp_dao_v2/models"
	"stp_dao_v2/utils"
	"strings"
	"time"
)

type Service struct {
	serverName string
	serverMark string
	gConfig    *oo.Config
	appConfig  *config.AppConfig
	scanInfo   []*config.ScanInfoConfig
	mCache     *cache.Cache
}

func NewService() *Service {
	return &Service{}
}

func (svc *Service) Start(ctx *cli.Context) error {
	defer func() {
		if err := recover(); err != nil {
			oo.LogW("got panic: %v", err)
		}
	}()

	svc.initServerName(ctx)
	svc.initMainLogger()

	if err := svc.loadGlobalConfig(ctx); err != nil {
		oo.LogW("load global config failed: %v", err)
		return err
	}
	if err := svc.initGlobalMysqlPool(); err != nil {
		oo.LogW("init mysql failed: %v", err)
		return err
	}
	svc.mCache = cache.New(cache.NoExpiration, time.Duration(24)*time.Hour)
	// STPT holder, snapshot at #15781071
	svc.mCache.Set(consts.CacheTokenHolders, 13804, cache.NoExpiration)

	go svc.scheduledTask()
	go svc.updateDaoInfoTask()
	go tokensImgTask()
	go updateNotification()
	go updateAccountRecord()
	go svc.getV1Proposal()
	go svc.getEthTokenHoldersTotal()

	router := gin.Default()
	router.Use(utils.Cors())

	basePath := svc.appConfig.BasePath

	r1 := router.Group(path.Join(basePath, "/img"))
	{
		r1.POST("/upload", svc.httpUploadImg)
	}
	r2 := router.Group(path.Join(basePath, "/dao"))
	{
		r2.GET("/list", svc.httpDaoList)
		r2.POST("/member", httpDaoJoinOrQuit)
		r2.GET("/left", httpLeftDaoJoin)
		r2.GET("/info", httpDaoInfo)
		r2.GET("/admins", httpDaoAdmins)
	}
	r3 := router.Group(path.Join(basePath, "/proposal"))
	{
		r3.GET("/list", httpProposalsList)
		r3.POST("/save", httpSaveProposal)
		r3.GET("/query", httpQueryProposal)
		r3.GET("/snapshot", httpQuerySnapshot)
	}
	r4 := router.Group(path.Join(basePath, "/votes"))
	{
		r4.GET("/list", httpVotesList)
	}
	r5 := router.Group(path.Join(basePath, "/sign"))
	{
		r5.POST("/create", svc.httpCreateSign)
		r5.POST("/lock/handle", svc.httpLockDaoHandleSign)
		r5.GET("/query/handle", svc.httpQueryDaoHandle)
	}
	r6 := router.Group(path.Join(basePath, "/account"))
	{
		r6.POST("/query", httpQueryAccount)
		r6.POST("/update", httpUpdateAccount)
		r6.GET("/record", httpQueryRecordList)
		r6.GET("/sign/list", httpQueryAccountSignList)
	}
	r7 := router.Group(path.Join(basePath, "/token"))
	{
		r7.GET("/list", httpTokenList)
		r7.GET("/img", httpTokenImg)
	}
	r8 := router.Group(path.Join(basePath, "/error"))
	{
		r8.POST("/info", httpErrorInfo)
	}
	r9 := router.Group(path.Join(basePath, "/airdrop"))
	{
		r9.POST("/create", svc.httpCreateAirdrop)
		r9.GET("/collect", httpCollectInformation)
		r9.POST("/save/user", httpSaveUserInformation)
		r9.POST("/user/download", httpDownloadUserInformation)
		r9.POST("/address", httpSaveAirdropAddress)
		r9.GET("/address/list", httpAirdropAddressList)
		r9.GET("/proof", httpClaimAirdrop)
	}
	r10 := router.Group(path.Join(basePath, "/activity"))
	{
		r10.GET("/list", httpActivity)
	}
	r11 := router.Group(path.Join(basePath, "/notification"))
	{
		r11.GET("/list", httpNotificationList)
		r11.POST("/read", httpNotificationRead)
		r11.GET("/unread/total", httpNotificationUnreadTotal)
	}
	r12 := router.Group(path.Join(basePath, "/overview"))
	{
		r12.GET("/total", svc.httpRecordTotal)
	}

	url := ginSwagger.URL(svc.appConfig.SwaggerUrl)
	router.GET(path.Join(basePath, "/swagger/*any"), ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router.Run(fmt.Sprintf("0.0.0.0:%d", svc.appConfig.ServerPort))
}

func (svc *Service) initServerName(ctx *cli.Context) {
	svc.serverName = strings.Split(filepath.Base(os.Args[0]), ".")[0]
	svc.serverMark = oo.GetSvrmark(svc.serverName)
}

func (svc *Service) initMainLogger() {
	oo.InitLog("./", svc.serverName, svc.serverName, func(str string) {})
}

func (svc *Service) initGlobalMysqlPool() error {
	var (
		conf config.MysqlConfig
		err  error
	)
	if err = svc.gConfig.SessDecode(svc.appConfig.MysqlConfName, &conf); err != nil {
		return err
	}
	oo.GMysqlPool, err = oo.InitMysqlPool(conf.Host, (int32)(conf.Port), conf.User, conf.Password, conf.Name)
	return err
}

func (svc *Service) loadGlobalConfig(ctx *cli.Context) error {
	configFile := ctx.String("config")
	initDomain := ctx.String("init-domain")

	var err error
	// config
	svc.gConfig, err = oo.InitConfig(configFile, nil)
	if err != nil {
		return err
	}
	if err = svc.gConfig.SessDecode(initDomain, &svc.appConfig); err != nil {
		return err
	}

	size := len(svc.appConfig.ScanInfoConfName)
	if size > 0 {
		for index := 0; index < size; index++ {
			var scanInfo *config.ScanInfoConfig
			if err = svc.gConfig.SessDecode(svc.appConfig.ScanInfoConfName[index], &scanInfo); err != nil {
				return err
			}
			svc.scanInfo = append(svc.scanInfo, scanInfo)
		}
	}

	return nil
}

func checkLogin(sign *models.SignData) (ret bool) {
	ret, errSign := utils.CheckPersonalSign(consts.SignMessagePrefix, sign.Account, sign.Signature)
	if errSign != nil {
		oo.LogD("signMessage err %v", errSign)
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
	oo.LogD("count:%v", count)
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

func (svc *Service) getEthTokenHoldersTotal() {
	defer time.AfterFunc(time.Duration(10*60)*time.Second, svc.getEthTokenHoldersTotal)

	var entities []models.DaoModel
	sqlSel := oo.NewSqler().Table(consts.TbNameDao).Select()
	err := oo.SqlSelect(sqlSel, &entities)
	if err != nil {
		oo.LogW("SQL err:%v", err)
		return
	}

	var queryTokenAddress []string
	var set = make(map[string]bool)
	for index := range entities {
		if set[entities[index].TokenAddress] && set[string(entities[index].TokenChainId)] {
			continue
		}
		set[entities[index].TokenAddress] = true
		set[string(entities[index].TokenChainId)] = true

		if entities[index].TokenChainId == 1 {
			queryTokenAddress = append(queryTokenAddress, entities[index].TokenAddress)
		}
	}

	var totalTokenHolder uint64
	for _, token := range queryTokenAddress {

		for indexScan := range svc.scanInfo {
			for indexUrl := range svc.scanInfo[indexScan].ChainId {
				if svc.scanInfo[indexScan].ChainId[indexUrl] == 1 {
					url := svc.scanInfo[indexScan].ScanUrl[indexUrl]

					res, err := utils.QueryErc20TokenHolders(token, url)
					if err != nil {
						oo.LogW("QueryErc20TokenHolders err:%v", err)
						return
					}
					totalTokenHolder += res.Result.Total

					key := fmt.Sprintf(`%d-%s`, 1, token)
					svc.mCache.Set(key, res.Result.Total, cache.NoExpiration)

				}
			}
		}
	}

	svc.mCache.Set(consts.CacheTokenHolders, totalTokenHolder, cache.NoExpiration)
}
