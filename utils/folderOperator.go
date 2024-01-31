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

type FolderOperator struct {
}

type DirectoryContent struct {
	Name     string
	IsDir    bool
	Mode     string
	ModeTime string
	Size     int64
	Path     string
	children []*DirectoryContent
}

func folderOperatorInstance() *FolderOperator {
	single := FolderOperator{}
	return &single
}

// CreateDirectory mkdir
func (fo *FolderOperator) CreateDirectory(dirname string) error {
	var absPath = dirname
	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), dirname)
	}

	err := os.MkdirAll(absPath, os.ModePerm)
	return err
}

// DeleteDirectory rm -rf
func (fo *FolderOperator) DeleteDirectory(dirname string) error {
	var absPath = dirname
	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), dirname)
	}

	if isDir := Common.IsDir(absPath); !isDir {
		return fmt.Errorf("%s is not a directory", absPath)
	}

	err := os.RemoveAll(absPath)
	return err
}

// Tree :获取目标文件夹下及子文件夹下的所有内容 类似于Linux中的tree命令
func (fo *FolderOperator) Tree(dirname string) (string, error) {
	dirname = Common.CompletePath(dirname)
	fmt.Println("dirname", dirname)
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
			dir, _ := circulateDirectory(SplicingPath(path, entry.Name(), string(os.PathSeparator)), depth)
			s = SplicingPath(s, strings.Repeat("----", depth-1), entry.Name(), "\n", dir)
		} else {
			s = SplicingPath(s, strings.Repeat("----", depth-1), entry.Name(), "\n")
		}
		depth--
	}

	return s, nil
}

// ListDirectoryContents :获取目标文件夹下的所有内容
func (fo *FolderOperator) ListDirectoryContents(dirname string) ([]*DirectoryContent, error) {
	var absPath = dirname
	//contentList := make([]consts.StrDict, 0, 10)
	contents := make([]*DirectoryContent, 0)

	if isAbs := filepath.IsAbs(absPath); !isAbs {
		absPath = SplicingPath(os.Getenv("workDirectory"), string(os.PathSeparator), dirname)
	}

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

// Ls : ls
func (fo *FolderOperator) Ls(dirname string) (string, error) {
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

// Lsl : ls -l
func (fo *FolderOperator) Lsl(dirname string) (string, error) {
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
func (fo *FolderOperator) GetAllFolderDepthContent(dirname string) ([]*DirectoryContent, error) {
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
				children: directoryContents,
			})
		} else {
			contents = append(contents, &DirectoryContent{
				Name:     entry.Name(),
				IsDir:    false,
				Mode:     fileInfo.Mode().String(),
				ModeTime: fileInfo.ModTime().Format("Jan 02 2006 15:04:05"),
				Size:     fileInfo.Size(),
				Path:     absPath,
				children: nil,
			})
		}
	}

	return contents, nil
}

func (fo *FolderOperator) CopyDirectory(src, dst string) error {
	// sourcePath      : D:\\awesomeGolang\\src\\base\\consts\\a
	// destinationPath : C:\\Users\\v_hhdzhang\\Downloads\\consts\\a
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

	for _, content := range contents {
		srcPath := filepath.Join(src, content.Name)
		dstPath := filepath.Join(dst, content.Name)
		if content.IsDir {
			err := fo.CopyDirectory(srcPath, dstPath)
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
	return nil
}

func (fo *FolderOperator) Cpr(srcDirname, dstDirname string) error {
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

	err = fo.CopyDirectory(srcPath, dstPath)
	if err != nil {
		return err
	}

	return nil
}
