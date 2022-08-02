package main

import (
	"fmt"
	"github.com/urfave/cli"
	"log"
	"net/http"
	"os"
	"stp_dao_v2/dao/controller"
)

func main() {
	if os.Getenv("debugPProf") == "true" {
		go func() {
			// terminal: $ go tool pprof -http=:8081 http://localhost:6060/debug/pprof/heap
			// web:
			// 1、http://localhost:8081/ui
			// 2、http://localhost:6060/debug/charts
			// 3、http://localhost:6060/debug/pprof
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	app := cli.NewApp()

	app.Name = "stp_dao_v2"
	app.Version = "v0.1.0"
	app.Description = "stp_dao_v2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config,c",
			Usage: "config file",
			Value: "./config/config.toml",
		},
		cli.StringFlag{
			Name:  "init-domain,i",
			Usage: "toml config file init domain",
			Value: "stp_dao_v2",
		},
	}
	server := controller.NewService()
	app.Action = server.Start

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
