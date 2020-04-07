package histogram

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type histogram map[rune]uint64

var l sync.Mutex

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

	h := histogram{}
	ch := checkFiles(files)
	for value := range ch {
		l.Lock()
		h[value]++
		l.Unlock()
	}

	return &h, nil
}

func checkFiles(files []string) chan rune {
	resCh := make(chan rune)

	go func() {
		for _, filePath := range files {

			for fileRes := range checkFile(filePath) {
				resCh <- fileRes
			}

		}
		close(resCh)
	}()

	return resCh
}

func checkFile(filePath string) chan rune {
	resCh := make(chan rune)

	go func() {
		file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}

		reader := bufio.NewReader(file)
		for err == nil {
			var b byte
			b, err = reader.ReadByte()
			if err == io.EOF {
				continue
			}
			if err != nil {
				panic(err)
				break
			}
			resCh <- rune(b)
		}

		close(resCh)
	}()

	return resCh
}
