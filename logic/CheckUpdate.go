package logic

import (
	"gf2gacha/util"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func CheckUpdate() (string, error) {
	resp, err := http.Get("https://gfl2worker.mcc.wiki/gf2gacha/version")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New(resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	newVersion := string(bodyBytes)
	if util.GetVersion() != newVersion {
		return newVersion, nil
	}
	return "", nil
}
