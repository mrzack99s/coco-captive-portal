package installer_utils

import (
	"io/fs"
	"os/exec"
)

type CommandType struct {
	Type    string
	Name    string
	Command exec.Cmd
}

type DownloadType struct {
	Name            string
	URL             string
	DestinationFile string
}

type CopyType struct {
	Src  string
	Dst  string
	Perm fs.FileMode
}

type ReplaceWordInFileType struct {
	OldWord string
	NewWord string
	File    string
}
