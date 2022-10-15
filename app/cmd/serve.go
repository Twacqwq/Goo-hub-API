package cmd

import (
	"thub/bootstrap"
	"thub/pkg/config"
	"thub/pkg/console"
	"thub/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var CmdServe = &cobra.Command{
	Use:   "serve",
	Short: "Start web server",
	Run:   runWeb,
	Args:  cobra.NoArgs, // 限制无参 输入其他参数报错
}

func init() {
	rootCmd.AddCommand(CmdServe)
}

func runWeb(cmd *cobra.Command, args []string) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	bootstrap.SetupRoute(router)
	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		logger.ErrorString("CMD", "serve", err.Error())
		console.Exit("Unable to start server, error:" + err.Error())
	}
}
