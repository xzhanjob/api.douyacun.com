package helper

import (
	"errors"
	"testing"
)

func Test_Image_EncodeWebP(t *testing.T) {
	filePath := "/Users/liuning/Documents/gitlab/api.douyacun.com/storage/images/1/golang/assert/go-function.png"
	if err := Image.EncodeWebP(filePath); err != nil {
		t.Fatal(err)
	}
	expactFile := "/Users/liuning/Documents/gitlab/api.douyacun.com/storage/images/1/golang/assert/go-function.webp"
	if !File.IsFile(expactFile) {
		t.Fatal(errors.New("webp生成失败"))
	}
}
