package main

import (
	"fmt"
	"log"
	"time"

	"github.com/242617/other/tasks/histogram/histogram"
)

func main() {
	start := time.Now()
	res, err := histogram.Scan("./files")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(`done in %f, found %d chars`+"\n", time.Since(start).Seconds(), len(*res))

	for char, count := range *res {
		fmt.Println(string(char), count)
	}
}
