package main

import (
	"log"
	"os"

	"github.com/242617/other/zip/zip"
)

func main() {
	log.SetFlags(log.Lshortfile)

	checkArgs(1, "not enough arguments")

	switch os.Args[1] {

	case "zip":
		checkArgs(4, "not enough arguments for zip")
		target := os.Args[2]
		files := os.Args[3:]
		err := zip.Zip(target, files...)
		if err != nil {
			log.Fatal("err", err)
		}

	case "unzip":
		checkArgs(4, "not enough arguments for unzip")
		source := os.Args[2]
		target := os.Args[3]
		err := zip.UnZip(source, target)
		if err != nil {
			log.Fatal("err", err)
		}

	default:
		log.Fatal("unknown mode")

	}

}

func checkArgs(n int, s string) {
	if len(os.Args) < n {
		log.Fatal(s)
	}
}
