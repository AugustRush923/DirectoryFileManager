package utils

import (
	"base/consts"
	"fmt"
	"os"
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
			fmt.Println("读取文件错误：", err)
			return err
		}

		if exist {
			return fmt.Errorf("文件%s已存在", filename)
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
		fmt.Println("读取文件错误：", err)
		return err
	}
	if !exist {
		return fmt.Errorf("文件%s不存在", filename)
	}

	if isDir := Common.IsFile(absPath); !isDir {
		return fmt.Errorf("%s is not a file", absPath)
	}

	err = os.WriteFile(absPath, []byte(content), 0666)

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
		fmt.Println("读取文件错误：", err)
		return consts.EmptyString, err
	}
	if !exist {
		return consts.EmptyString, fmt.Errorf("文件%s不存在", filename)
	}

	fileContent, err := os.ReadFile(absPath)

	return string(fileContent), err
}

// DeleteFile :删除指定位置的文件
func (f *fileOperator) DeleteFile(filename string) error {
	var absPath = Common.CompleteFullPath(filename)

	exist, err := Common.IsExist(absPath)
	if err != nil {
		fmt.Println("读取文件错误：", err)
		return err
	}
	if !exist {
		return fmt.Errorf("文件%s不存在", filename)
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

	err = f.CreateFile(newAbsPath, true)
	if err != nil {
		return err
	}

	err = f.WriteFile(newAbsPath, content)
	if err != nil {
		return err
	}

	return nil
}
