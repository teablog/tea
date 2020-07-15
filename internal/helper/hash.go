package helper

import (
	"crypto/md5"
	"encoding/hex"
)

func Md532(buf []byte) string {
	h := md5.New()
	h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}

func Md516(buf []byte) string {
	return Md532(buf)[8:24]
}
