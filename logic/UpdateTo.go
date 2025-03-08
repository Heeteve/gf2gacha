package logic

import (
	"fmt"
	"github.com/inconshreveable/go-update"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/exec"
)

func UpdateTo(version string) error {
	resp, err := http.Get(fmt.Sprintf("https://gfl2bucket.mcc.wiki/gf2gacha/%s/gf2gacha.exe", version))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		return err
	}

	err = restart()
	if err != nil {
		return err
	}
	os.Exit(0)

	return nil
}

func restart() error {
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath)
	return cmd.Start()
}
