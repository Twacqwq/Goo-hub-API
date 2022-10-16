package make

import (
	"fmt"
	"strings"
	"thub/pkg/console"

	"github.com/spf13/cobra"
)

var CmdMakeAPIController = &cobra.Command{
	Use:   "apicontroller",
	Short: "Create api controller, example: make apicontroller v1/user",
	Run:   runMakeAPIController,
	Args:  cobra.ExactArgs(1),
}

func runMakeAPIController(cmd *cobra.Command, arg []string) {
	array := strings.Split(arg[0], "/")
	if len(array) != 2 {
		console.Exit("api controller name format: v1/user")
	}

	// apiVersion 用来拼接目标路径
	// name 用来生成 cmd.Model 实例
	apiVersion, name := array[0], array[1]
	model := makeModelFromString(name)

	filePath := fmt.Sprintf("app/http/controllers/api/%s/%s_controller.go", apiVersion, model.TableName)

	createFileFromStub(filePath, "apicontroller", model)
}
