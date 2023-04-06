package routers

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"stp_dao_v2/handlers"
	"stp_dao_v2/middlewares"

	_ "stp_dao_v2/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	basePath  string
	GinEngine *gin.Engine
)

// Init sets GinEngine as custom gin engine and uses custom middlewares
func init() {
	// Init a gin engine
	GinEngine = gin.New()

	// Add error middleware
	GinEngine.Use(middlewares.ErrHandler())

	// Add preset cors middleware, you can add your own middleware
	GinEngine.Use(middlewares.Cors())
	// Add default recovery middleware, you can add your own middleware
	// GinEngine.Use(gin.Recovery())
}

func Router() {
	// Setup the base uri
	basePath = viper.GetString("app.base_path")
	router := GinEngine.Group(path.Join(basePath))
	{
		r1 := router.Group(path.Join(basePath, "/img"))
		{
			r1.POST("/upload", handlers.HttpUploadImg)
		}
		r2 := router.Group(path.Join(basePath, "/dao"))
		{
			r2.GET("/list", handlers.HttpDaoList)
			r2.POST("/member", handlers.HttpDaoJoinOrQuit)
			r2.GET("/left", handlers.HttpLeftDaoJoin)
			r2.GET("/info", handlers.HttpDaoInfo)
			r2.GET("/admins", handlers.HttpDaoAdmins)
		}
		r3 := router.Group(path.Join(basePath, "/proposal"))
		{
			r3.GET("/list", handlers.HttpProposalsList)
			r3.POST("/save", handlers.HttpSaveProposal)
			r3.GET("/query", handlers.HttpQueryProposal)
			r3.GET("/snapshot", handlers.HttpQuerySnapshot)
		}
		r4 := router.Group(path.Join(basePath, "/votes"))
		{
			r4.GET("/list", handlers.HttpVotesList)
		}
		r5 := router.Group(path.Join(basePath, "/sign"))
		{
			r5.POST("/create", handlers.HttpCreateSign)
			r5.POST("/lock/handle", handlers.HttpLockDaoHandleSign)
			r5.GET("/query/handle", handlers.HttpQueryDaoHandle)
		}
		r6 := router.Group(path.Join(basePath, "/account"))
		{
			r6.POST("/query", handlers.HttpQueryAccount)
			r6.POST("/update", handlers.HttpUpdateAccount)
			r6.GET("/record", handlers.HttpQueryRecordList)
			r6.GET("/sign/list", handlers.HttpQueryAccountSignList)
			r6.GET("/nfts", handlers.HttpQueryAccountNFTsList)
			r6.POST("/update/follow", handlers.HttpUpdateAccountFollow)
			r6.GET("/following/list", handlers.HttpAccountFollowingList)
			r6.GET("/followers/list", handlers.HttpAccountFollowersList)
			r6.GET("/relation", handlers.HttpAccountRelation)
			r6.GET("/push/setting", handlers.HttpPushSetting)
		}
		r7 := router.Group(path.Join(basePath, "/token"))
		{
			r7.GET("/list", handlers.HttpTokenList)
			r7.GET("/img", handlers.HttpTokenImg)
		}
		r8 := router.Group(path.Join(basePath, "/error"))
		{
			r8.POST("/info", handlers.HttpErrorInfo)
		}
		r9 := router.Group(path.Join(basePath, "/airdrop"))
		{
			r9.POST("/create", handlers.HttpCreateAirdrop)
			r9.GET("/collect", handlers.HttpCollectInformation)
			r9.POST("/save/user", handlers.HttpSaveUserInformation)
			r9.POST("/user/download", handlers.HttpDownloadUserInformation)
			r9.POST("/address", handlers.HttpSaveAirdropAddress)
			r9.GET("/address/list", handlers.HttpAirdropAddressList)
			r9.GET("/proof", handlers.HttpClaimAirdrop)
		}
		r10 := router.Group(path.Join(basePath, "/activity"))
		{
			r10.GET("/list", handlers.HttpActivity)
		}
		r11 := router.Group(path.Join(basePath, "/notification"))
		{
			r11.GET("/list", handlers.HttpNotificationList)
			r11.POST("/read", handlers.HttpNotificationRead)
			r11.GET("/unread/total", handlers.HttpNotificationUnreadTotal)
		}
		r12 := router.Group(path.Join(basePath, "/overview"))
		{
			r12.GET("/total", handlers.HttpRecordTotal)
		}
		r13 := router.Group(path.Join(basePath, "/swap"))
		{
			r13.POST("/create", handlers.CreateSwap)
			r13.POST("/purchased", handlers.PurchasedSwap)
			r13.GET("/list", handlers.SwapList)
			r13.GET("/transactions", handlers.SwapTransactionsList)
			r13.GET("/prices", handlers.SwapPrices)
			r13.GET("/isWhite", handlers.SwapIsWhite)
			r13.GET("/isCreatorWhite", handlers.SwapIsCreatorWhite)
		}
	}

	url := ginSwagger.URL(viper.GetString("app.swagger_url"))
	GinEngine.GET(path.Join(basePath, "/swagger/*any"), ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	GinEngine.Run(fmt.Sprintf(":%d", viper.GetInt("app.server_port")))
}
