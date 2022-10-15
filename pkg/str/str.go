// Package str 字符串辅助方法
package str

import (
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

// Plural 转为复数 user -> users
func Plural(word string) string {
	return pluralize.NewClient().Plural(word)
}

// Signgular 转为单数 users -> user
func Signgular(word string) string {
	return pluralize.NewClient().Singular(word)
}

// 驼峰转下划线命名
func Snake(s string) string {
	return strcase.ToSnake(s)
}

// 下划线转驼峰命名
func Camel(s string) string {
	return strcase.ToCamel(s)
}

// 大驼峰转小驼峰命名法
func LowerCamel(s string) string {
	return strcase.ToLowerCamel(s)
}
