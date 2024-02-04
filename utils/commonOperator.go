package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

type commonOperator struct{}

var Common = commonOperatorInstance()

func init() {
	Common.SetWorkDirectory()
}

func commonOperatorInstance() *commonOperator {
	// 单例
	single := &commonOperator{}
	return single
}

// ShowWorkDirectory pwd
func (co *commonOperator) ShowWorkDirectory() string {
	return os.Getenv("workDirectory")
}

// SetWorkDirectory 设置/变更当前工作目录
func (co *commonOperator) SetWorkDirectory(paths ...string) bool {
	var path string
	if len(paths) > 0 {
		path = paths[0]
	}

	// 如果没有接收到变量则默认为当前所在目录
	if path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			fmt.Println("get current work directory failed: err=", err)
			return false
		}
		path = currentPath
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Println("change directory failed: err=", err)
		return false
	}

	err = os.Setenv("workDirectory", path)
	if err != nil {
		fmt.Println("set environment failed: err=", err)
		return false
	}
	return true
}

// IsDir 判断给出路径是否为目录
func (co *commonOperator) IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("get path: %s info failed: err=%v", path, err)
		return false
	}
	return fileInfo.IsDir()
}

// IsFile 判断给出路径是否为文件
func (co *commonOperator) IsFile(path string) bool {
	return !co.IsDir(path)
}

// IsExist 判断给出路径是否存在
func (co *commonOperator) IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RenameOrMove 重命名/移动指定目录/文件
func (co *commonOperator) RenameOrMove(oldFilename, newFilename string) error {
	var oldAbsPath = oldFilename
	if isAbs := filepath.IsAbs(oldAbsPath); !isAbs {
		oldAbsPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), oldFilename)
	}
	var newAbsPath = newFilename
	if isAbs := filepath.IsAbs(newAbsPath); !isAbs {
		newAbsPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), newFilename)
	}

	err := os.Rename(oldAbsPath, newAbsPath)

	return err
}

// CompleteFullPath 完善给到路径为绝对路径
func (co *commonOperator) CompleteFullPath(path string) string {
	if isAbs := filepath.IsAbs(path); !isAbs {
		path = filepath.Join(os.Getenv("workDirectory"), path)
	}

	return path
}
