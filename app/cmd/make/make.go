package make

import (
	"embed"
	"fmt"
	"strings"
	"thub/pkg/console"
	"thub/pkg/file"
	"thub/pkg/str"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

// Model 参数解释
// 单个词，用户命令传参，以 User 模型为例：
//   - user
//   - User
//   - users
//   - Users
//
// 整理好的数据：
//
//	{
//	    "TableName": "users",
//	    "StructName": "User",
//	    "StructNamePlural": "Users"
//	    "VariableName": "user",
//	    "VariableNamePlural": "users",
//	    "PackageName": "user"
//	}
//
// -
// 两个词或者以上，用户命令传参，以 TopicComment 模型为例：
//   - topic_comment
//   - topic_comments
//   - TopicComment
//   - TopicComments
//
// 整理好的数据：
//
//	{
//	    "TableName": "topic_comments",
//	    "StructName": "TopicComment",
//	    "StructNamePlural": "TopicComments"
//	    "VariableName": "topicComment",
//	    "VariableNamePlural": "topicComments",
//	    "PackageName": "topic_comment"
//	}
type Model struct {
	TableName          string
	StructName         string
	StructNamePlural   string
	VariableName       string
	VariableNamePlural string
	PackageName        string
}

//go:embed stubs
var stubsFS embed.FS

var CmdMake = &cobra.Command{
	Use:   "make",
	Short: "Generate file and code",
}

func init() {
	CmdMake.AddCommand(
		CmdMakeCMD,
	)
}

// 格式化用户输入的内容
func makeModelFromString(name string) Model {
	model := Model{}
	model.StructName = str.Signgular(strcase.ToCamel(name))
	model.StructNamePlural = str.Plural(model.StructName)
	model.TableName = str.Snake(model.StructNamePlural)
	model.VariableName = str.LowerCamel(model.StructName)
	model.PackageName = str.Snake(model.StructName)
	model.VariableNamePlural = str.LowerCamel(model.StructNamePlural)

	return model
}

// 读取 stub 文件并进行变量替换
// 最后一个选项可选 传参应传 map[string]string 类型
func createFileFromStub(filePath, stubName string, model Model, variables ...map[string]string) {
	replaces := make(map[string]string)
	if len(variables) > 0 {
		replaces = variables[0]
	}

	// 目标文件已存在
	if file.Exists(filePath) {
		console.Exit(filePath + " already exists!")
	}

	// 读取模板文件
	modelData, err := stubsFS.ReadFile("stubs/" + stubName + ".stub")
	if err != nil {
		console.Exit(err.Error())
	}
	modelStub := string(modelData)

	// 添加默认的替换变量
	replaces["{{VariableName}}"] = model.VariableName
	replaces["{{VariableNamePlural}}"] = model.VariableNamePlural
	replaces["{{StructName}}"] = model.StructName
	replaces["{{StructNamePlural}}"] = model.StructNamePlural
	replaces["{{TableName}}"] = model.TableName
	replaces["{{PackageName}}"] = model.PackageName

	// 对模板内容进行替换
	for search, replace := range replaces {
		modelStub = strings.ReplaceAll(modelStub, search, replace)
	}

	err = file.Put([]byte(modelStub), filePath)
	if err != nil {
		console.Exit(err.Error())
	}

	console.Success(fmt.Sprintf("[%v] created.", filePath))
}
