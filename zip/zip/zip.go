package zip

import (
	"archive/zip"
	"io"
	"log"
	"os"
)

func Zip(target string, files ...string) error {

	targetFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("err", err)
		return err
	}

	w := zip.NewWriter(targetFile)

	for _, filename := range files {

		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		writer, err := w.Create(filename)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = targetFile.Close()
	if err != nil {
		return err
	}

	return nil

}
