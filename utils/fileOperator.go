package utils

import (
	"base/consts"
	"fmt"
	"os"
	"path/filepath"
)

var File = fileOperatorInstance()

type fileOperator struct {
}

func fileOperatorInstance() *fileOperator {
	single := fileOperator{}
	return &single
}

// CreateFile :创建文件
func (f *fileOperator) CreateFile(filename string, override bool) error {
	var absPath = filename

	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), filename)
	}

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

// WriteFile echo
func (f *fileOperator) WriteFile(filename, content string) error {
	var absPath = filename
	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), filename)
	}

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

// ReadFile cat/more/less/tail/head
func (f *fileOperator) ReadFile(filename string) (string, error) {
	var absPath = filename
	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), filename)
	}

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

// DeleteFile rm
func (f *fileOperator) DeleteFile(filename string) error {
	var absPath = filename
	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), filename)
	}

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

// CopyFile cp
func (f *fileOperator) CopyFile(oldFileName, newFileName string) error {
	content, err := f.ReadFile(oldFileName)
	if err != nil {
		return err
	}

	err = f.CreateFile(newFileName, true)
	if err != nil {
		return err
	}

	err = f.WriteFile(newFileName, content)
	if err != nil {
		return err
	}

	return nil
}
