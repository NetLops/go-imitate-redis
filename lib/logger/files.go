package logger

import (
	"fmt"
	"os"
)

// checkNotExist
//
// Description: 检查文件是否出现
//
//
// param: src string
//
// return: bool
func checkNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

// checkPermission
//
// Description: 检查是否有响应权限
//
//
// param: src string
//
// return: bool
func checkPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// mkDir
//
// Description: 创建目录（包含嵌套级别目录）
//
//
// param: src string
//
// return: error
func mkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// isNotExistMkdir
//
// Description: 判断文件没有就创建
//
//
// param: src string
//
// return: error
func isNotExistMkdir(src string) error {
	if notExist := checkNotExist(src); notExist == true {
		if err := mkDir(src); err != nil {
			return err
		}
	}
	return nil
}

// mustOpen
//
// Description: 打开文件
//
//
// param: fileName string
// param: dir string
//
// return: *os.File
// return: error
func mustOpen(fileName, dir string) (*os.File, error) {
	perm := checkPermission(dir)
	if perm == true {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	err := isNotExistMkdir(dir)
	if err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %s", dir, err)
	}

	file, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}
	return file, nil
}
