package cmd

import (
	"fmt"
	"os"

	appctx "github.com/hoangtk0100/app-context"
	"github.com/hoangtk0100/app-context/component/datastore/gormdb"
	"github.com/hoangtk0100/app-context/component/dbmigration"
	ginserver "github.com/hoangtk0100/app-context/component/server/gin"
	"github.com/hoangtk0100/app-context/component/token"
	"github.com/hoangtk0100/stauvise/pkg/common"
	"github.com/hoangtk0100/stauvise/pkg/handler"
	"github.com/spf13/cobra"
)

func newAppContext() appctx.AppContext {
	return appctx.NewAppContext(
		appctx.WithName("stauvise"),
		appctx.WithComponent(gormdb.NewGormDB(common.PluginDBMain, common.PluginDBMain)),
		appctx.WithComponent(token.NewPasetoMaker(common.PluginTokenMaker)),
		appctx.WithComponent(dbmigration.NewDBMigration(common.PluginDBMigrator)),
		appctx.WithComponent(ginserver.NewServer(common.PluginGin)),
	)
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "Start STAUVISE service",
	Run: func(cmd *cobra.Command, args []string) {
		appCtx := newAppContext()
		log := appCtx.Logger("service")

		if err := appCtx.Load(); err != nil {
			log.Fatal(err)
		}

		config := common.NewConfig(appCtx)
		server := handler.NewServer(config)

		server.Start()
	},
}

func Execute() {
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
