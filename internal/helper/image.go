package helper

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chai2010/webp"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var Image _Image

type _Image struct{}

// 图片转webp, 会在相应的图片未知生成相应的webp
func (*_Image) EncodeWebP(filepath string) (err error) {
	var (
		img image.Image
	)
	if Image.IsWebP(filepath) {
		return errors.New("文件已经是webp")
	}
	var buf bytes.Buffer
	//file, err := os.Open("/Users/liuning/Documents/linux/www/foo/asserts/images/groutine-cover.jpeg")
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	ext := path.Ext(filepath)
	switch ext {
	case ".png":
		img, err = png.Decode(file)
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	default:
		return errors.New(fmt.Sprintf("暂不支持%s图片转WebP", ext))
	}
	if err != nil {
		return err
	}
	file.Close()
	dst := fmt.Sprintf("%s/%s.webp", path.Dir(filepath), strings.TrimRight(path.Base(filepath), ext))
	if err = webp.Encode(&buf, img, &webp.Options{Lossless: false, Quality: 10}); err != nil {
		log.Println(err)
	}
	if err = ioutil.WriteFile(dst, buf.Bytes(), 0666); err != nil {
		log.Println(err)
	}
	return
}

func (*_Image) IsWebP(filePath string) bool {
	ext := path.Ext(filePath)
	return ext == ".webp"
}

func (*_Image) WebPSupportExt(ext string) bool {
	switch ext {
	case ".png", ".jpg", ".jpeg":
		return true
	default:
		return false
	}
}
