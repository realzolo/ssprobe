package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func PathRewrite2(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if _, err := os.Stat(abs); err != nil && os.IsNotExist(err) {
		abs = strings.Replace(filepath.ToSlash(abs), "server/", "", 1)
	}
	return abs
}
func PathRewrite(path string) string {
	abs, _ := filepath.Abs("")
	path = filepath.Clean(path)

	fakePath := filepath.Join(abs, path)
	if _, err := os.Stat(fakePath); err != nil && os.IsNotExist(err) {
		path = filepath.ToSlash(path)
		fakePath = filepath.ToSlash(fakePath)

		if strings.Index(path, "/") == 0 {
			path = strings.TrimLeft(path, "/")
		}
		if strings.LastIndex(fakePath, "/") == 0 {
			fakePath = strings.TrimRight(path, "/")
		}
		return strings.Replace(fakePath, path, "server/"+path, -1)
	}
	return fakePath
}
