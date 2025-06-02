package logic

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func AppendLog(accessToken string, uid int64, gachaUrl string) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return errors.WithStack(err)
	}

	logPath := filepath.Join(userHome, "/AppData/LocalLow/SunBorn/少女前线2：追放/Player.log")
	file, err := os.OpenFile(logPath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf(`{"access_token":"%s","uid":%d}\n`, accessToken, uid))
	if err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf(`{"gacha_record_url":"%s"}\n`, gachaUrl))
	if err != nil {
		return err
	}

	return nil
}
