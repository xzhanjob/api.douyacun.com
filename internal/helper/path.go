package helper

import (
	"bytes"
	"strings"
)

var Path _path

type _path struct{}

// 拼接多个路径，最后不包含 /
// Example:
// Join("/ab/c", "/c", "f/e")
// /ab/c/c/f/e
func (*_path) Join(opts ...string) string {
	path := bytes.NewBufferString("")
	path.WriteString("/")
	for _, v := range opts {
		path.WriteString(strings.Trim(v, "/"))
		path.WriteString("/")
	}
	p := strings.TrimRight(path.String(), "/")
	return p
}

// 获取上层目录
func (*_path) Dir(filePath string) string {
	index := strings.LastIndex(filePath, "/")
	if index == -1 {
		return filePath
	}
	return filePath[:index]
}

// 获取当前文件名
func (*_path) File(filePath string) string {
	index := strings.LastIndex(filePath, "/")
	if index == -1 {
		return filePath
	}
	return filePath[index+1:]
}

// 判断是否为绝对路径
func (*_path) IsAbs(file string) bool {
	return strings.HasPrefix(file, "/")
}
