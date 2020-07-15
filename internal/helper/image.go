package helper

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
	"sync"
)

var Image _Image

type _Image struct{}

func (*_Image) Convert(dir string) error {
	wg := sync.WaitGroup{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir failed, %s", err)
	}
	for _, v := range files {
		if v.Name() == ".." || v.Name() == "." {
			continue
		}
		if v.Mode().IsDir() {
			_ = Image.Convert(path.Join(dir, v.Name()))
		} else if v.Mode().IsRegular() {
			if Image.WebPSupportExt(path.Ext(v.Name())) {
				wg.Add(1)
				go func(filePath string, wg *sync.WaitGroup) {
					defer wg.Done()
					_ = Image.EncodeWebP(filePath)
				}(path.Join(dir, v.Name()), &wg)
			}
		}
	}
	wg.Wait()
	return nil
}

// 图片转webp, 会在相应的图片未知生成相应的webp
func (*_Image)  EncodeWebP(filepath string) (err error) {
	if Image.IsWebP(filepath) {
		return errors.New("文件已经是webp")
	}
	ext := path.Ext(filepath)
	dst := strings.ReplaceAll(filepath, ext, ".webp")
	if _, err = exec.LookPath("cwebp"); err != nil {
		return fmt.Errorf("libwebp还没有安装")
	}
	cmd := exec.Command("cwebp", "-q", "10", filepath, "-o", dst)
	if err = cmd.Run(); err != nil {
		return
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
