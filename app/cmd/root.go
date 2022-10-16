package cmd

import (
	"fmt"
	"os"
	"thub/app/cmd/make"
	"thub/bootstrap"
	"thub/pkg/config"
	"thub/pkg/console"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Thub",
	Short: "A simple forum project",
	Long:  `Default will run "serve" command, you can use "-h" flag to see all subcommands`,

	// rootCmd 的所有子命令都会执行的代码 Hooks
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.InitConfig(Env)
		bootstrap.SetupLogger()
		bootstrap.SetupDB()
		bootstrap.SetupRedis()
	},
}

func init() {
	rootCmd.AddCommand(CmdServe)
	rootCmd.AddCommand(CmdKey)
	rootCmd.AddCommand(CmdPlay)
	rootCmd.AddCommand(make.CmdMake)
}

// 执行主命令
func Execute() {
	// 配置默认运行 web 服务
	RegisterDefaultCmd(rootCmd, CmdServe)
	// 注册全局参数, --env
	RegisterGlobalFlags(rootCmd)
	// 启动
	if err := rootCmd.Execute(); err != nil {
		console.Exit(fmt.Sprintf("Failed to run app with %v: %s", os.Args, err.Error()))
	}
}
