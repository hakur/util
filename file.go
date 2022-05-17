package util

import (
	"io"
	"os"
)

// FileGetContents read file into a string variable
// FileGetContents 把一个文件读成一个字符串 http://php.net/manual/en/function.file-get-contents.php
func FileGetContents(filename string) (content string) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	var buf = make([]byte, 256)
	for {
		n, err := f.Read(buf)
		if err != nil {
			return content
		}
		content += string(buf[0:n])
	}
}

// FileExists check file exists
// FileExists 查询一个文件是否存在
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsEmptyDir check path if is a empty directory
// IsEmptyDir 查询一个目录是否是空目录
func IsEmptyDir(dirPath string) bool {
	f, err := os.Open(dirPath)
	if err != nil {
		return false
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF { // Or f.Readdir(1)
		return true
	}
	return false
}

// ReadDirRecusive read files under directory with recursive
// ReadDirRecusive 递归读取目录下所有的文件和文件夹
func ReadDirRecusive(path string, callback func(filePath string) (err error)) (err error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range files {
		filepath := path + "/" + fi.Name()
		if callback != nil {
			if err = callback(filepath); err != nil {
				return err
			}
		}

		if fi.IsDir() {
			ReadDirRecusive(filepath, callback)
		}
	}

	return nil
}
