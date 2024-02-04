package utils

import (
	"base/consts"
	"base/models"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var Folder = folderOperatorInstance()

type folderOperator struct {
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
			s = SplicingString(s, strings.Repeat("----", depth-1), entry.Name(), "\n", dir)
		} else {
			s = SplicingString(s, strings.Repeat("----", depth-1), entry.Name(), "\n")
		}
		depth--
	}

	return s, nil
}

// ListDirectoryContents :获取目标文件夹下的所有内容 类似于Linux中的ls命令
func (fo *folderOperator) ListDirectoryContents(dirname string) ([]*models.DirectoryContent, error) {
	contents := make([]*models.DirectoryContent, 0)

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
			contents = append(contents, &models.DirectoryContent{
				FileInfo: &models.FileInfo{
					Name:     contentInfo.Name(),
					IsDir:    contentInfo.IsDir(),
					Mode:     contentInfo.Mode().String(),
					ModeTime: contentInfo.ModTime().Format("Jan 02 2006 15:04:05"),
					Size:     contentInfo.Size(),
				},
			})
		}
	} else {
		contents = append(contents, &models.DirectoryContent{
			FileInfo: &models.FileInfo{
				Name:     dirname,
				IsDir:    false,
				Mode:     fileInfo.Mode().String(),
				ModeTime: fileInfo.ModTime().Format("Jan 02 2006 15:04:05"),
				Size:     fileInfo.Size(),
			},
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
		resultStr = SplicingString(resultStr, content.Name, "\t")
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
		resultStr = SplicingString(resultStr, content.Mode, "\t", strconv.FormatInt(content.Size, 10), "\t", content.ModeTime, "\t", content.Name, "\n")
	}

	return resultStr, nil
}

// GetAllFolderDepthContent :递归获取文件夹下的所有内容
func (fo *folderOperator) GetAllFolderDepthContent(dirname string) ([]*models.DirectoryContent, error) {
	absPath := dirname
	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingString(os.Getenv("workDirectory"), string(os.PathSeparator), dirname)
	}
	contents := make([]*models.DirectoryContent, 0)

	isDir := Common.IsDir(absPath)
	if !isDir {
		return nil, fmt.Errorf("%s is not a directory", absPath)
	}

	if !strings.HasSuffix(absPath, string(os.PathSeparator)) {
		absPath = SplicingString(absPath, string(os.PathSeparator))
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
			directoryContents, _ := fo.GetAllFolderDepthContent(SplicingString(absPath, entry.Name(), string(os.PathSeparator)))
			contents = append(contents, &models.DirectoryContent{
				FileInfo: &models.FileInfo{
					Name:     entry.Name(),
					IsDir:    true,
					Mode:     fileInfo.Mode().String(),
					ModeTime: fileInfo.ModTime().Format("Jan 02 2006 15:04:05"),
					Size:     fileInfo.Size(),
				},
				Children: directoryContents,
			})
		} else {
			contents = append(contents, &models.DirectoryContent{
				FileInfo: &models.FileInfo{
					Name:     entry.Name(),
					IsDir:    false,
					Mode:     fileInfo.Mode().String(),
					ModeTime: fileInfo.ModTime().Format("Jan 02 2006 15:04:05"),
					Size:     fileInfo.Size(),
				},
				Children: nil,
			})
		}
	}

	return contents, nil
}

// CopyDirectory :拷贝目录下的所有内容至新路径下 类似于Linux中的cp -r命令
func (fo *folderOperator) CopyDirectory(srcDirname, dstDirname string) error {
	srcPath := Common.CompleteFullPath(srcDirname)
	srcPathList := strings.Split(srcPath, string(os.PathSeparator))
	dirname := srcPathList[len(srcPathList)-1]

	dstPath := Common.CompleteFullPath(dstDirname)
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

// RenameDirectory :重命名目录
func (fo *folderOperator) RenameDirectory(oldName, newName string) error {
	var oldAbsPath, newAbsPath = Common.CompleteFullPath(oldName), Common.CompleteFullPath(newName)
	var oldPathList, newPathList = strings.Split(oldAbsPath, string(os.PathSeparator)), strings.Split(newAbsPath, string(os.PathSeparator))

	if filepath.Join(oldPathList[:len(oldPathList)-1]...) !=
		filepath.Join(newPathList[:len(newPathList)-1]...) {
		return fmt.Errorf("%s and %s aren't at the same path", oldAbsPath, newAbsPath)
	}

	return Common.RenameOrMove(oldAbsPath, newAbsPath)
}

// MoveDirectory :移动目录至指定目录下
func (fo *folderOperator) MoveDirectory(oldName, newName string) error {
	var oldPath, newPath = Common.CompleteFullPath(oldName), Common.CompleteFullPath(newName)
	var oldPathList, newPathList = strings.Split(oldPath, string(os.PathSeparator)), strings.Split(newPath, string(os.PathSeparator))

	if filepath.Join(oldPathList[:len(oldPathList)-1]...) ==
		filepath.Join(newPathList[:len(newPathList)-1]...) {
		return fmt.Errorf("%s and %s are at the same path", oldPath, newPath)
	}

	// 拷贝至新目录下
	err := fo.CopyDirectory(oldPath, newPath)
	if err != nil {
		return err
	}

	// 删除源目录
	err = fo.DeleteDirectory(oldPath)
	if err != nil {
		return err
	}

	return nil
}
