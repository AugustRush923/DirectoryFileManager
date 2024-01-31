package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type CommonOperator struct{}

var Common = commonOperatorInstance()

func init() {
	Common.SetWorkDirectory()
}

func commonOperatorInstance() *CommonOperator {
	// 单例
	single := &CommonOperator{}
	return single
}

// ShowWorkDirectory pwd
func (co *CommonOperator) ShowWorkDirectory() string {
	return os.Getenv("workDirectory")
}

// SetWorkDirectory cd
func (co *CommonOperator) SetWorkDirectory(paths ...string) bool {
	// 设置/变更当前工作目录
	var path string
	if len(paths) > 0 {
		path = paths[0]
	}

	if path == "" {
		currentPath, err := os.Getwd()
		if err != nil {
			fmt.Println("获取当前工作目录失败")
			return false
		}
		path = currentPath
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Println("改变目录失败")
		return false
	}

	err = os.Setenv("workDirectory", path)
	if err != nil {
		fmt.Println("设置环境变量失败")
		return false
	}
	return true
}

func (co *CommonOperator) IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func (co *CommonOperator) IsFile(path string) bool {
	return !co.IsDir(path)
}

func (co *CommonOperator) IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// RenameOrMove mv
func (co *CommonOperator) RenameOrMove(oldFilename, newFilename string) error {
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

func (co *CommonOperator) CompletePath(path string) string {
	// 输入路径是绝对路径
	if isAbs := filepath.IsAbs(path); isAbs {
		// 路径是文件夹
		if isDir := co.IsDir(path); isDir {
			return SplicingPath(path, string(os.PathSeparator))
		}
		return path
	}
	// 输入路径是相对路径
	path = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), path)
	if isDir := co.IsDir(path); isDir && !strings.HasSuffix(path, string(os.PathSeparator)) {
		return SplicingPath(path, string(os.PathSeparator))
	}
	return path
}
