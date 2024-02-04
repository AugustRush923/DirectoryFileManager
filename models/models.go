package models

type FileInfo struct {
	Name     string
	IsDir    bool
	Mode     string
	ModeTime string
	Size     int64
}

type DirectoryContent struct {
	Info     *FileInfo
	Children []*DirectoryContent
}
