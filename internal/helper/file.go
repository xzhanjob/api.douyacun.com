package helper

import (
	"fmt"
	"io"
	"os"
	"path"
)

var File _File

type _File struct{}

// 判断是否为常规文件文件，文件是否存在
func (*_File) IsFile(filename string) bool {
	info, err := os.Stat(filename)

	return err == nil && info.Mode().IsRegular()
}

// Overwrite overwrites the file with data. Creates file if not present.
func (*_File) FileOverwrite(fileName string, data []byte) bool {
	f, err := os.Create(fileName)
	if err != nil {
		return false
	}

	_, err = f.Write(data)
	return err == nil
}

func (*_File) Copy(dst, src string) (int64, error) {
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()
	if err = os.MkdirAll(path.Dir(dst), 0755); err != nil {
		return 0, fmt.Errorf("目录 %s 创建失败 %s", path.Dir(dst), err)
	}
	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
// 判断目录是否存在，是否为目录
func (*_File) IsDir(dir string) bool {
	fi, err := os.Stat(dir)
	if err != nil {
		return false
	}
	if fi.Mode().IsDir() {
		return true
	}
	return false
}
