package helper

import (
	"errors"
	"testing"
)

func Test_Image_EncodeWebP(t *testing.T) {
	filePath := "/Users/liuning/Documents/gitlab/api.douyacun.com/storage/images/1/assert/go-内存分配-三阶段分配.png"
	if err := Image.EncodeWebP(filePath); err != nil {
		t.Fatal(err)
	}
	expactFile := "/Users/liuning/Documents/gitlab/api.douyacun.com/storage/images/1/assert/go-内存分配-三阶段分配.webp"
	if !FileExists(expactFile) {
		t.Fatal(errors.New("webp生成失败"))
	}
}
