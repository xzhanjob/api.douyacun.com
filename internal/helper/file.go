package helper

import (
	"fmt"
	"io"
	"os"
	"path"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)

	return err == nil && !info.IsDir()
}

// Overwrite overwrites the file with data. Creates file if not present.
func FileOverwrite(fileName string, data []byte) bool {
	f, err := os.Create(fileName)
	if err != nil {
		return false
	}

	_, err = f.Write(data)
	return err == nil
}

func Copy(dst, src string) (int64, error)  {
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
