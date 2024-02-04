package utils

import (
	"base/consts"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var File = fileOperatorInstance()

type fileOperator struct {
}

func fileOperatorInstance() *fileOperator {
	single := fileOperator{}
	return &single
}

// CreateFile :在指定位置创建文件 类似于Linux中的touch命令
func (f *fileOperator) CreateFile(filename string, override bool) error {
	absPath := Common.CompleteFullPath(filename)

	if !override {
		exist, err := Common.IsExist(absPath)
		if err != nil {
			fmt.Printf("get file %s failed：err=%v", absPath, err)
			return err
		}

		if exist {
			return fmt.Errorf("%s is already exist", absPath)
		}
	}

	file, err := os.Create(absPath)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

// WriteFile :往指定文件写入内容
func (f *fileOperator) WriteFile(filename, content string) error {
	var absPath = Common.CompleteFullPath(filename)

	exist, err := Common.IsExist(absPath)
	if err != nil {
		fmt.Printf("get file %s failed：err=%v", absPath, err)
		return err
	}
	if !exist {
		return fmt.Errorf("%s is not exist", absPath)
	}

	if isDir := Common.IsFile(absPath); !isDir {
		return fmt.Errorf("%s is not a file", absPath)
	}

	err = os.WriteFile(absPath, []byte(content), os.ModePerm)

	return err
}

// ReadFile :读取指定位置的文件全部内容
func (f *fileOperator) ReadFile(filename string) (string, error) {
	var absPath = Common.CompleteFullPath(filename)

	// 更简便的方式读取文件的所有内容
	if isDir := Common.IsFile(absPath); !isDir {
		return consts.EmptyString, fmt.Errorf("%s is not a file", absPath)
	}

	exist, err := Common.IsExist(absPath)
	if err != nil {
		fmt.Printf("get file %s failed：err=%v", absPath, err)
		return consts.EmptyString, err
	}
	if !exist {
		return consts.EmptyString, fmt.Errorf("%s is not exist", filename)
	}

	fileContent, err := os.ReadFile(absPath)

	return string(fileContent), err
}

// DeleteFile :删除指定位置的文件
func (f *fileOperator) DeleteFile(filename string) error {
	var absPath = Common.CompleteFullPath(filename)

	exist, err := Common.IsExist(absPath)
	if err != nil {
		fmt.Printf("get file %s failed：err=%v", absPath, err)
		return err
	}
	if !exist {
		return fmt.Errorf("%s is not exist", filename)
	}

	if isDir := Common.IsFile(absPath); !isDir {
		return fmt.Errorf("%s is not a file", absPath)
	}

	err = os.Remove(absPath)
	return err
}

// CopyFile :拷贝文件
func (f *fileOperator) CopyFile(oldFileName, newFileName string) error {
	var oldAbsPath, newAbsPath = Common.CompleteFullPath(oldFileName), Common.CompleteFullPath(newFileName)

	content, err := f.ReadFile(oldAbsPath)
	if err != nil {
		return err
	}

	err = f.CreateFile(newAbsPath, false)
	if err != nil {
		return err
	}

	err = f.WriteFile(newAbsPath, content)
	if err != nil {
		return err
	}

	return nil
}

// RenameFile :重命名文件
func (f *fileOperator) RenameFile(oldFileName, newFileName string) error {
	var oldPath, newPath = Common.CompleteFullPath(oldFileName), Common.CompleteFullPath(newFileName)
	var oldPathList, newPathList = strings.Split(oldPath, string(os.PathSeparator)), strings.Split(newPath, string(os.PathSeparator))

	if filepath.Join(oldPathList[:len(oldPathList)-1]...) !=
		filepath.Join(newPathList[:len(newPathList)-1]...) {
		return fmt.Errorf("%s and %s are not at the same path", oldPath, newPath)
	}

	return Common.RenameOrMove(oldPath, newPath)
}

// MoveFile :移动文件到指定位置
func (f *fileOperator) MoveFile(oldFileName, newFileName string) error {
	var oldPath, newPath = Common.CompleteFullPath(oldFileName), Common.CompleteFullPath(newFileName)
	var oldPathList, newPathList = strings.Split(oldPath, string(os.PathSeparator)), strings.Split(newPath, string(os.PathSeparator))

	if filepath.Join(oldPathList[:len(oldPathList)-1]...) ==
		filepath.Join(newPathList[:len(newPathList)-1]...) {
		return fmt.Errorf("%s and %s are at the same path", oldPath, newPath)
	}

	return Common.RenameOrMove(oldPath, newPath)
}
