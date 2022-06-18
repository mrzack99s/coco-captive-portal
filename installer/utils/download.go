package installer_utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func downloadPackages() (err error) {
	packages := []DownloadType{}
	log.Info().Msg("# download coco-captive-portal packages")

	packages = append(packages, DownloadType{
		Name:            "coco-captive-portal",
		URL:             "https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/coco",
		DestinationFile: fmt.Sprintf("%s/coco", APP_DIR),
	})

	packages = append(packages, DownloadType{
		Name:            "coco-captive-portal service",
		URL:             "https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/coco-captive-portal.service",
		DestinationFile: "/etc/systemd/system/coco-captive-portal.service",
	})

	packages = append(packages, DownloadType{
		Name:            "sample config",
		URL:             "https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/config.yaml.sample",
		DestinationFile: fmt.Sprintf("%s/config.yaml.sample", APP_DIR),
	})

	packages = append(packages, DownloadType{
		Name:            "coco-dist-ui",
		URL:             "https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/dist-ui.tar.gz",
		DestinationFile: "/tmp/coco-dist-ui.tar.gz",
	})

	for _, dl := range packages {
		log.Info().Msg(getDownloadMessage(dl, DOING_STATE))
		if e := getDownloadAndSave(dl.URL, dl.DestinationFile); e != nil {
			if IGNORE_VERIFY {
				log.Warn().Msg(getDownloadMessage(dl, FAILED_STATE))
			} else {
				log.Error().Msg(getDownloadMessage(dl, FAILED_STATE))
				err = e
				return
			}
		} else {
			log.Info().Msg(getDownloadMessage(dl, DONE_STATE))
		}
	}
	return
}

func getDownloadAndSave(link, location string) (err error) {

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Put content on file
	resp, err := client.Get(link)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Create blank file
	file, err := os.Create(location)
	if err != nil {
		return
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return
	}
	defer file.Close()

	return
}
