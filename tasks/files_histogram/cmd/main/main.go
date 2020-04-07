package main

import (
	"fmt"
	"log"
	"time"

	"github.com/242617/other/tasks/files_histogram/histogram"
)

func main() {

	start := time.Now()
	res, err := histogram.Check("./folder")
	fmt.Printf(`done in %f, found %d chars`+"\n", time.Since(start).Seconds(), len(*res))
	die(err)

	for char, count := range *res {
		fmt.Println(string(char), count)
	}
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
