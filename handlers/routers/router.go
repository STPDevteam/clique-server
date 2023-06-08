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
		r1 := router.Group("/img")
		{
			r1.POST("/upload", handlers.HttpUploadImg)
		}
		r2 := router.Group("/dao")
		{
			r2.GET("/list", handlers.HttpDaoList)
			r2.POST("/member", handlers.HttpDaoJoinOrQuit)
			r2.GET("/left", handlers.HttpLeftDaoJoin)
			r2.GET("/info", handlers.HttpDaoInfo)
			r2.GET("/admins", handlers.HttpDaoAdmins)
			r2.GET("/one", handlers.HttpDaoOne)
		}
		r3 := router.Group("/proposal")
		{
			r3.GET("/list", handlers.HttpProposalsList)
			r3.POST("/save", handlers.HttpSaveProposal)
			r3.GET("/query", handlers.HttpQueryProposal)
			r3.GET("/snapshot", handlers.HttpQuerySnapshot)
		}
		r4 := router.Group("/votes")
		{
			r4.GET("/list", handlers.HttpVotesList)
		}
		r5 := router.Group("/sign")
		{
			r5.POST("/create", handlers.HttpCreateSign)
			r5.POST("/lock/handle", handlers.HttpLockDaoHandleSign)
			r5.GET("/query/handle", handlers.HttpQueryDaoHandle)
		}
		account := router.Group("/account")
		{
			account.POST("/jwt/signIn", handlers.HttpAccountSignIn)
			account.GET("/record", handlers.HttpQueryRecordList)
			account.GET("/sign/list", handlers.HttpQueryAccountSignList)
			account.GET("/nfts", handlers.HttpQueryAccountNFTsList)
			account.GET("/following/list", handlers.HttpAccountFollowingList)
			account.GET("/followers/list", handlers.HttpAccountFollowersList)
			account.GET("/relation", handlers.HttpAccountRelation)
			account.GET("/top/list", handlers.HttpAccountTopList)
		}
		accountAuthForce := router.Group("/account", middlewares.JWTAuthForce())
		{
			accountAuthForce.POST("/update", handlers.HttpUpdateAccount)
			accountAuthForce.POST("/update/follow", handlers.HttpUpdateAccountFollow)
			accountAuthForce.POST("/push/setting", handlers.HttpPushSetting)
		}
		accountAuth := router.Group("/account", middlewares.JWTAuth())
		{
			accountAuth.POST("/query", handlers.HttpQueryAccount)
		}
		r7 := router.Group("/token")
		{
			r7.GET("/list", handlers.HttpTokenList)
			r7.GET("/img", handlers.HttpTokenImg)
		}
		r8 := router.Group("/error")
		{
			r8.POST("/info", handlers.HttpErrorInfo)
		}
		r9 := router.Group("/airdrop")
		{
			r9.POST("/create", handlers.HttpCreateAirdrop)
			r9.GET("/collect", handlers.HttpCollectInformation)
			r9.POST("/save/user", handlers.HttpSaveUserInformation)
			r9.POST("/user/download", handlers.HttpDownloadUserInformation)
			r9.POST("/address", handlers.HttpSaveAirdropAddress)
			r9.GET("/address/list", handlers.HttpAirdropAddressList)
			r9.GET("/proof", handlers.HttpClaimAirdrop)
		}
		r10 := router.Group("/activity")
		{
			r10.GET("/list", handlers.HttpActivity)
		}
		r11 := router.Group("/notification")
		{
			r11.GET("/list", handlers.HttpNotificationList)
			r11.POST("/read", handlers.HttpNotificationRead)
			r11.GET("/unread/total", handlers.HttpNotificationUnreadTotal)
		}
		r12 := router.Group("/overview")
		{
			r12.GET("/total", handlers.HttpRecordTotal)
		}
		r13 := router.Group("/swap")
		{
			r13.POST("/create", handlers.CreateSwap)
			r13.POST("/purchased", handlers.PurchasedSwap)
			r13.GET("/list", handlers.SwapList)
			r13.GET("/transactions", handlers.SwapTransactionsList)
			r13.GET("/prices", handlers.SwapPrices)
			r13.GET("/isWhite", handlers.SwapIsWhite)
			r13.GET("/isCreatorWhite", handlers.SwapIsCreatorWhite)
		}

		jobs := router.Group("/jobs")
		{
			jobs.POST("/publish", middlewares.JWTAuthForce(), handlers.JobsPublish)
			jobs.POST("/publish/edit", middlewares.JWTAuthForce(), handlers.JobsPublishEdit)
			jobs.GET("/publish/list", handlers.JobsPublishList)
			jobs.DELETE("/publish/:jobPublishId", middlewares.JWTAuthForce(), handlers.JobsPublishDelete)

			jobs.POST("/apply", middlewares.JWTAuthForce(), handlers.JobsApply)
			jobs.POST("/join/member", middlewares.JWTAuthForce(), handlers.JobsJoinToMember)
			jobs.POST("/apply/review", middlewares.JWTAuthForce(), handlers.JobsApplyReview)
			jobs.POST("/alter", middlewares.JWTAuthForce(), handlers.JobsAlter)
			jobs.GET("/identity", middlewares.JWTAuthForce(), handlers.JobsIdentity)
			jobs.GET("/left", middlewares.JWTAuthForce(), handlers.JobsLeft)

			jobs.GET("/apply/list", handlers.JobsApplyList)
			jobs.GET("/list", handlers.JobsList)
		}

		teamSpaces := router.Group("/spaces")
		{
			teamSpaces.POST("/create", middlewares.JWTAuthForce(), handlers.CreateTeamSpaces)
			teamSpaces.POST("/update", middlewares.JWTAuthForce(), handlers.UpdateTeamSpaces)
			teamSpaces.POST("/remove", middlewares.JWTAuthForce(), handlers.TeamSpacesRemoveToTrash)
			teamSpaces.POST("/delete", middlewares.JWTAuthForce(), handlers.DeleteTeamSpaces)

			teamSpaces.GET("/list", middlewares.JWTAuth(), handlers.TeamSpacesList)
		}

		task := router.Group("/task")
		{
			task.POST("/create", middlewares.JWTAuthForce(), handlers.CreateTask)
			task.POST("/update", middlewares.JWTAuthForce(), handlers.UpdateTask)
			task.POST("/remove", middlewares.JWTAuthForce(), handlers.TaskRemoveToTrash)
			task.POST("/delete", middlewares.JWTAuthForce(), handlers.DeleteTask)

			task.GET("/list", handlers.TaskList)
			task.GET("/detail/:taskId", handlers.TaskDetail)
		}

		sbt := router.Group("/sbt")
		{
			sbt.POST("/create", middlewares.JWTAuthForce(), handlers.CreateSBT)

			sbt.GET("/list", handlers.SBTList)
		}

		ai := router.Group("/ai")
		{
			ai.POST("", middlewares.JWTAuthForce(), handlers.Ai)
		}
	}

	url := ginSwagger.URL(viper.GetString("app.swagger_url"))
	GinEngine.GET(path.Join(basePath, "/swagger/*any"), ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	GinEngine.Run(fmt.Sprintf(":%d", viper.GetInt("app.server_port")))
}
