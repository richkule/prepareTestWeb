package helpFun

import (
	"crypto/md5"
	"encoding/hex"
	"io"
)

// Возвращает hex-строку md5 хэша передаваемой строки
func HexMD5(text string) (string, error) {
	h := md5.New()
	if _, err := io.WriteString(h, text); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
