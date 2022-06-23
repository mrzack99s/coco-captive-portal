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

	dlLink := []string{
		"https://github.com/mrzack99s/coco-captive-portal/releases/download/" + APP_VERSION + "/coco",
		"https://github.com/mrzack99s/coco-captive-portal/releases/download/" + APP_VERSION + "/coco-captive-portal.service",
		"https://github.com/mrzack99s/coco-captive-portal/releases/download/" + APP_VERSION + "/config.yaml.sample",
		"https://github.com/mrzack99s/coco-captive-portal/releases/download/" + APP_VERSION + "/dist-auth-ui.tar.gz",
		"https://github.com/mrzack99s/coco-captive-portal/releases/download/" + APP_VERSION + "/dist-operator-ui.tar.gz",
	}
	if LATEST {
		dlLink = []string{
			"https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/coco",
			"https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/coco-captive-portal.service",
			"https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/config.yaml.sample",
			"https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/dist-auth-ui.tar.gz",
			"https://github.com/mrzack99s/coco-captive-portal/releases/latest/download/dist-operator-ui.tar.gz",
		}
	}

	packages = append(packages, DownloadType{
		Name:            "coco-captive-portal",
		URL:             dlLink[0],
		DestinationFile: fmt.Sprintf("%s/coco", APP_DIR),
	})

	packages = append(packages, DownloadType{
		Name:            "coco-captive-portal service",
		URL:             dlLink[1],
		DestinationFile: "/etc/systemd/system/coco-captive-portal.service",
	})

	packages = append(packages, DownloadType{
		Name:            "sample config",
		URL:             dlLink[2],
		DestinationFile: fmt.Sprintf("%s/config.yaml.sample", APP_DIR),
	})

	packages = append(packages, DownloadType{
		Name:            "coco-dist-auth-ui",
		URL:             dlLink[3],
		DestinationFile: "/tmp/coco-dist-auth-ui.tar.gz",
	})

	packages = append(packages, DownloadType{
		Name:            "coco-dist-operator-ui",
		URL:             dlLink[4],
		DestinationFile: "/tmp/coco-dist-operator-ui.tar.gz",
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
