package fs

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// md5gen 實際改用傻2，因爲 md5 已經被認爲不安全， lint 工具一直報警
func md5gen(fullPath string) (string, error) {
	data, err := ioutil.ReadFile(filepath.Clean(fullPath))
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, err = h.Write(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// MD5Gen generate a md5sum file for file "a" with filename "a.md5"
func MD5Gen(fullPath string) (bool, error) {
	md5, err := md5gen(fullPath)
	if err != nil {
		return false, err
	}
	f, err := os.Create(fullPath + ".md5")
	if err != nil {
		return false, err
	}
	defer f.Close()
	if _, err = f.WriteString(md5); err != nil {
		return false, err
	}
	if err = f.Sync(); err != nil {
		return false, err
	}

	return true, nil
}

// CheckMD5 verify the file's md5sum, confirm it's not modified.
func CheckMD5(fullPath string) bool {
	md5, err := md5gen(fullPath)
	if err != nil {
		return false
	}

	fi, err := os.Open(filepath.Clean(fullPath + ".md5"))
	if err != nil {
		return false
	}
	defer fi.Close()

	buf := make([]byte, 64)
	n, err := fi.Read(buf)
	if n != 64 || err != nil {
		return false
	}
	if strings.Compare(md5, string(buf)) != 0 {
		return false
	}
	return true
}
