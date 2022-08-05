package installer_utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

func ReplaceInFile(t ReplaceWordInFileType) (err error) {
	log.Info().Msgf("opening %s file", t.File)
	input, err := ioutil.ReadFile(t.File)
	if err != nil {
		return
	}
	log.Info().Msgf("replacing %s to %s in %s file", t.OldWord, t.NewWord, t.File)
	output := bytes.Replace(input, []byte(t.OldWord), []byte(t.NewWord), -1)
	err = ioutil.WriteFile(t.File, output, 0644)
	if err == nil {
		log.Info().Msgf("%s to %s in %s file was replaced", t.OldWord, t.NewWord, t.File)
	}
	return
}

func AppendStringToFile(t AppendStringToFileType) (err error) {
	log.Info().Msgf("opening %s file", t.File)
	f, err := os.OpenFile(t.File,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", t.Str)); err == nil {
		log.Info().Msgf("added %s to %s file", t.Str, t.File)
	}
	return
}
