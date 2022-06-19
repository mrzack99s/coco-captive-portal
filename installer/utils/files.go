package installer_utils

import (
	"bytes"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

func Replace(t ReplaceWordInFileType) (err error) {
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
