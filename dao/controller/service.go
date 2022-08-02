package controller

import (
	_ "crypto/ecdsa"
	"fmt"
	oo "github.com/Anna2024/liboo"
	"github.com/gin-gonic/gin"
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
)

type Service struct {
	serverName string
	serverMark string
	gConfig    *oo.Config
	appConfig  *config.AppConfig
	scanInfo   []*config.ScanInfoConfig
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

	go svc.DoScheduledTask()
	go DoUpdateDaoInfoTask()

	router := gin.Default()
	router.Use(utils.Cors())

	basePath := svc.appConfig.BasePath

	r1 := router.Group(path.Join(basePath, "/img"))
	{
		r1.POST("/upload", svc.httpUploadImg)
		r1.GET("/show", httpShowImg)
	}
	r2 := router.Group(path.Join(basePath, "/dao"))
	{
		r2.GET("/list", httpDaoList)
		r2.POST("/member", httpDaoJoinOrQuit)
		r2.GET("/left", httpLeftDaoCreator)
	}
	r3 := router.Group(path.Join(basePath, "/proposal"))
	{
		r3.GET("/list", httpProposalsList)
	}
	r4 := router.Group(path.Join(basePath, "/votes"))
	{
		r4.GET("/list", httpVotesList)
	}
	r5 := router.Group(path.Join(basePath, "/sign"))
	{
		r5.POST("/create", svc.httpCreateSign)
	}
	r6 := router.Group(path.Join(basePath, "/account"))
	{
		r6.POST("/query", httpQueryAccount)
		r6.POST("/update", httpUpdateAccount)
	}
	r7 := router.Group(path.Join(basePath, "/token"))
	{
		r7.GET("/list", httpTokenList)
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
			if len(scanInfo.SupportedCoin) != len(scanInfo.SupportedCoinAddress) || len(scanInfo.SupportedCoin) != len(scanInfo.SupportedCoinDecimals) {
				return fmt.Errorf("supportedCoin length mismatch")
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
		return
	}
	if !ret {
		oo.LogD("check Sign fail")
		return
	}
	return
}
