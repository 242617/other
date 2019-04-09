package zip

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_All(t *testing.T) {

	var content = map[string][]byte{
		"file1.txt": []byte(strings.Repeat("file 1 sample content\n", 1000)),
		"file2.txt": []byte(strings.Repeat("file 2 sample content\n", 999)),
	}
	var files []string

	// Creating files
	for filename, data := range content {
		err := ioutil.WriteFile(filename, data, 0644)
		if err != nil {
			t.Fatal("err", err)
		}
		files = append(files, filename)
	}

	// Compress
	err := Zip("target.zip", files...)
	if err != nil {
		t.Fatal("err", err)
	}

	// Remove created files
	for filename := range content {
		err := os.Remove(filename)
		if err != nil {
			t.Fatal("err", err)
		}
	}

	// Uncompress
	err = UnZip("target.zip", ".")
	if err != nil {
		t.Fatal("err", err)
	}

	// Checking content of extracted files
	for filename, data := range content {
		barr, err := ioutil.ReadFile(filename)
		if err != nil {
			t.Fatal("err", err)
		}

		if bytes.Compare(barr, data) != 0 {
			t.Log("barr", string(barr))
			t.Log("data", string(data))
			t.Fatal("not equal content")
		}
	}

	// Remove extracted files
	for filename := range content {
		err := os.Remove(filename)
		if err != nil {
			t.Fatal("err", err)
		}
	}

	// Remove created archive
	err = os.Remove("target.zip")
	if err != nil {
		t.Fatal("err", err)
	}

}
