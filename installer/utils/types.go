package installer_utils

import "os/exec"

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

type ReplaceWordInFileType struct {
	OldWord string
	NewWord string
	File    string
}
