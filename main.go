package main

import (
	"flag"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"stp_dao_v2/config"
	"stp_dao_v2/handlers/routers"
	"stp_dao_v2/tasks"
	"time"

	oo "github.com/Anna2024/liboo"
)

func init() {
	flag.StringVar(&config.ConfigFile, "c", "./config/config.toml", "The path of configuration file.")
	flag.Parse()
}

func main() {
	var err error
	viper.SetConfigFile(config.ConfigFile)
	err = viper.ReadInConfig()
	if err != nil {
		oo.LogW("viper.ReadInConfig err %v", err)
		return
	}
	viper.WatchConfig()

	config.ServerMark = oo.GetSvrmark(config.ServerName)
	svrTag := config.ServerName + "." + config.GitServer
	oo.InitLog("./", config.ServerName, svrTag, func(str string) {})

	defer func() {
		if err := recover(); nil != err {
			oo.LogW("panic err: %v", err)
		}
	}()

	oo.GMysqlPool, err = oo.InitMysqlPool(
		viper.GetString("mysql.host"),
		viper.GetInt32("mysql.port"),
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.name"),
	)
	if err != nil {
		oo.LogW("Failed to init mysql pool: %v", err)
		return
	}

	config.MyCache = cache.New(cache.NoExpiration, time.Duration(24)*time.Hour)

	go tasks.ScheduledTask()
	go tasks.UpdateDaoInfoTask()
	go tasks.TokensImgTask()
	go tasks.UpdateNotification()
	go tasks.UpdateSBTStatus()
	//go tasks.UpdateAccountRecord()
	//go tasks.DaoCountTask()
	//go tasks.SwapTokenPrice()
	//go tasks.UpdateSwapStatus()
	//go tasks.DoPush()

	//go tasks.OnceTaskForTeamSpaces()
	//go tasks.GetV1Proposal()

	go oo.SafeGuardTask(func() { routers.Router() }, time.Second*10)

	oo.WaitExitSignal()
}
