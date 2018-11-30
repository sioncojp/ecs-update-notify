package pidfile

import (
	"errors"
	"os"
	"path/filepath"
)

func Create(pidfile string) error {
	if pidfile == "" {
		return errors.New("pidfile not configured")
	}
	if err := os.MkdirAll(filepath.Dir(pidfile), os.FileMode(0755)); err != nil {
		return err
	}

	file, err := os.Create(pidfile)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func Remove(pidfile string) error {
	if err := os.Remove(pidfile); err != nil {
		return err
	}
	return nil
}
