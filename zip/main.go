package main

import (
	"log"
	"os"

	"github.com/242617/other/zip/zip"
)

func main() {
	log.SetFlags(log.Lshortfile)

	if len(os.Args) < 4 {
		log.Fatal("not enough arguments")
	}

	switch os.Args[1] {

	case "zip":
		target := os.Args[2]
		files := os.Args[3:]
		err := zip.Zip(target, files...)
		if err != nil {
			log.Fatal("err", err)
		}

	case "unzip":
		source := os.Args[2]
		target := os.Args[3]
		err := zip.UnZip(source, target)
		if err != nil {
			log.Fatal("err", err)
		}

	}

}
