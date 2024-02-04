package utils

import (
	"base/consts"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Folder = folderOperatorInstance()

type folderOperator struct {
}

type DirectoryContent struct {
	Name     string
	IsDir    bool
	Mode     string
	ModeTime string
	Size     int64
	Path     string
	Children []*DirectoryContent
}

func folderOperatorInstance() *folderOperator {
	single := folderOperator{}
	return &single
}

// CreateDirectory :在指定位置创建目录 类似于Linux中的mkdir -p命令
func (fo *folderOperator) CreateDirectory(dirname string) error {
	err := os.MkdirAll(Common.CompleteFullPath(dirname), os.ModePerm)
	return err
}

// DeleteDirectory :删除指定目录及目录下的所有内容 类似于Linux中的rm -rf命令
func (fo *folderOperator) DeleteDirectory(dirname string) error {
	absPath := Common.CompleteFullPath(dirname)

	if isDir := Common.IsDir(absPath); !isDir {
		return fmt.Errorf("%s is not a directory", absPath)
	}

	err := os.RemoveAll(absPath)
	return err
}

// Tree :获取目标文件夹下及子文件夹下的所有内容 类似于Linux中的tree命令
func (fo *folderOperator) Tree(dirname string) (string, error) {
	dirname = Common.CompleteFullPath(dirname)
	return circulateDirectory(dirname, 0)
}

func circulateDirectory(path string, depth int) (string, error) {
	var s string

	isDir := Common.IsDir(path)
	if !isDir {
		return consts.EmptyString, fmt.Errorf("%s is not a directory", path)
	}

	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return consts.EmptyString, err
	}

	for _, entry := range dirEntry {
		depth++
		if entry.IsDir() {
			dir, _ := circulateDirectory(filepath.Join(path, entry.Name()), depth)
			s = SplicingPath(s, strings.Repeat("----", depth-1), entry.Name(), "\n", dir)
		} else {
			s = SplicingPath(s, strings.Repeat("----", depth-1), entry.Name(), "\n")
		}
		depth--
	}

	return s, nil
}

// ListDirectoryContents :获取目标文件夹下的所有内容 类似于Linux中的ls命令
func (fo *folderOperator) ListDirectoryContents(dirname string) ([]*DirectoryContent, error) {
	contents := make([]*DirectoryContent, 0)

	absPath := Common.CompleteFullPath(dirname)

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path %s can't find: %v", absPath, err)
	}

	if fileInfo.IsDir() {
		dirEntry, err := os.ReadDir(absPath)

		if err != nil {
			return nil, fmt.Errorf("path %s cant't read: %v", absPath, err)
		}

		for _, entry := range dirEntry {
			contentInfo, err := entry.Info()

			if err != nil {
				return nil, fmt.Errorf("entry can't read: %v", err)
			}
			contents = append(contents, &DirectoryContent{
				Name:     contentInfo.Name(),
				IsDir:    contentInfo.IsDir(),
				Mode:     contentInfo.Mode().String(),
				ModeTime: contentInfo.ModTime().Format("Jan 02 2006 15:04:05"),
				Size:     contentInfo.Size(),
			})
		}
	} else {
		contents = append(contents, &DirectoryContent{
			Name:     dirname,
			IsDir:    false,
			Mode:     fileInfo.Mode().String(),
			ModeTime: fileInfo.ModTime().Format("Jan 02 2006 15:04:05"),
			Size:     fileInfo.Size(),
		})
	}

	return contents, nil
}

// Ls :获取目标文件夹下的所有内容 以字符串形式展示
func (fo *folderOperator) Ls(dirname string) (string, error) {
	contents, err := fo.ListDirectoryContents(dirname)
	if err != nil {
		return "", err
	}

	resultStr := consts.EmptyString
	for _, content := range contents {
		resultStr = SplicingPath(resultStr, content.Name, "\t")
	}

	return resultStr, nil
}

// Lsl : 获取目标文件夹下的所有内容的详细信息 以字符串形式展示
func (fo *folderOperator) Lsl(dirname string) (string, error) {
	contents, err := fo.ListDirectoryContents(dirname)

	if err != nil {
		return "", err
	}

	resultStr := consts.EmptyString
	for _, content := range contents {
		resultStr = SplicingPath(resultStr, content.Mode, "\t", strconv.FormatInt(content.Size, 10), "\t", content.ModeTime, "\t", content.Name, "\n")
	}

	return resultStr, nil
}

// GetAllFolderDepthContent :递归获取文件夹下的所有内容
func (fo *folderOperator) GetAllFolderDepthContent(dirname string) ([]*DirectoryContent, error) {
	absPath := dirname
	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), dirname)
	}
	contents := make([]*DirectoryContent, 0)

	isDir := Common.IsDir(absPath)
	if !isDir {
		return nil, fmt.Errorf("%s is not a directory", absPath)
	}

	if !strings.HasSuffix(absPath, string(os.PathSeparator)) {
		absPath = SplicingPath(absPath, string(os.PathSeparator))
	}

	dirEntry, err := os.ReadDir(absPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntry {
		fileInfo, err := entry.Info()
		if err != nil {
			return nil, err
		}

		if entry.IsDir() {
			directoryContents, _ := fo.GetAllFolderDepthContent(SplicingPath(absPath, entry.Name(), string(os.PathSeparator)))
			contents = append(contents, &DirectoryContent{
				Name:     entry.Name(),
				IsDir:    true,
				Mode:     fileInfo.Mode().String(),
				ModeTime: fileInfo.ModTime().Format("Jan 02 2006 15:04:05"),
				Size:     fileInfo.Size(),
				Path:     absPath,
				Children: directoryContents,
			})
		} else {
			contents = append(contents, &DirectoryContent{
				Name:     entry.Name(),
				IsDir:    false,
				Mode:     fileInfo.Mode().String(),
				ModeTime: fileInfo.ModTime().Format("Jan 02 2006 15:04:05"),
				Size:     fileInfo.Size(),
				Path:     absPath,
				Children: nil,
			})
		}
	}

	return contents, nil
}

// CopyDirectory :拷贝目录下的所有内容至新路径下 类似于Linux中的cp -r命令
func (fo *folderOperator) CopyDirectory(srcDirname, dstDirname string) error {
	srcPath := filepath.Clean(srcDirname)
	if isAbs := filepath.IsAbs(srcPath); !isAbs {
		srcPath = filepath.Join(os.Getenv("workDirectory"), srcPath)
	}
	srcPathList := strings.Split(srcPath, string(os.PathSeparator))
	dirname := srcPathList[len(srcPathList)-1]

	dstPath := filepath.Clean(dstDirname)
	if isAbs := filepath.IsAbs(dstPath); !isAbs {
		srcPath = filepath.Join(os.Getenv("workDirectory"), dstPath)
	}
	dstPath = filepath.Join(dstPath, dirname)

	err := fo.CreateDirectory(dstPath)
	if err != nil {
		return fmt.Errorf("create directory %s failed: %v", dstPath, err)
	}

	err = fo.copyDirectory(srcPath, dstPath)
	if err != nil {
		return err
	}

	return nil
}

func (fo *folderOperator) copyDirectory(src, dst string) error {
	src, dst = filepath.Clean(src), filepath.Clean(dst)

	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcFileInfo.IsDir() {
		return fmt.Errorf("source %s is not a directory", src)
	}

	// 获取源目录下的所有内容
	contents, err := fo.ListDirectoryContents(src)
	if err != nil {
		return err
	}
	if len(contents) == 0 {
		err := fo.CreateDirectory(dst)
		if err != nil {
			return err
		}
	} else {
		for _, content := range contents {
			srcPath := filepath.Join(src, content.Name)
			dstPath := filepath.Join(dst, content.Name)
			if content.IsDir {
				err := fo.copyDirectory(srcPath, dstPath)
				if err != nil {
					return err
				}
			} else {
				err := fo.CreateDirectory(filepath.Dir(dstPath))
				if err != nil {
					return err
				}

				err = File.CopyFile(srcPath, dstPath)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
