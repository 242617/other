package histogram

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type histogram map[rune]uint64

func Check(folder string) (*histogram, error) {

	fileInfos, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		files = append(files, folder+"/"+fileInfo.Name())
	}

	resCh := make(chan *histogram)
	errCh := make(chan error)

	go func() {
		var l sync.Mutex
		h := histogram{}
		ch := checkFiles(files, errCh)
		for value := range ch {
			l.Lock()
			h[value]++
			l.Unlock()
		}
		resCh <- &h
	}()

	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return nil, err
	}
}

func checkFiles(files []string, errCh chan error) chan rune {
	filesCh := make(chan rune)

	go func() {
		for _, filePath := range files {

			for fileRes := range checkFile(filePath, errCh) {
				filesCh <- fileRes
			}

		}
		close(filesCh)
	}()

	return filesCh
}

func checkFile(filePath string, errCh chan error) chan rune {
	resCh := make(chan rune)

	go func() {
		file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			errCh <- err
		}
		defer func() {
			err = file.Close()
			if err != nil {
				errCh <- err
			}
		}()

		reader := bufio.NewReader(file)
		for err == nil {
			var b byte
			b, err = reader.ReadByte()
			if err == io.EOF {
				continue
			}
			if err != nil {
				errCh <- err
				break
			}
			resCh <- rune(b)
		}

		close(resCh)
	}()

	return resCh
}
