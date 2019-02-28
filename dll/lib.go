// +build windows,amd64
package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import "strings"

func main() {}

//export GetInt
func GetInt() int32 {
	return 42
}

//export SetInt
func SetInt(n int32) int32 {
	return n * n
}

//export GetString
func GetString() *C.char {
	return C.CString("hello!")
}

//export SetString
func SetString(str *C.char) *C.char {
	s := C.GoString(str)
	return C.CString(strings.ToUpper(s) + "-" + strings.ToLower(s) + "-" + "...")
}
