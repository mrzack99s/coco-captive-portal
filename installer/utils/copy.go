package installer_utils

import (
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

func copy(cp CopyType) (err error) {
	bytesRead, e := ioutil.ReadFile(cp.Src)
	if e != nil {
		err = e
		return
	}
	e = ioutil.WriteFile(cp.Dst, bytesRead, cp.Perm)
	if e != nil {
		err = e
		return
	}

	log.Info().Msgf("copied %s to %s", cp.Src, cp.Dst)
	return
}
