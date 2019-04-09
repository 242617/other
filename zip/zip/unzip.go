package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func UnZip(source, target string) error {

	r, err := zip.OpenReader(source)
	if err != nil {
		return err
	}

	for _, f := range r.File {

		path := filepath.Join(target, f.Name)
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(file, rc)
		if err != nil {
			return err
		}

		err = rc.Close()
		if err != nil {
			return err
		}

	}

	err = r.Close()
	if err != nil {
		return err
	}

	return nil

}
