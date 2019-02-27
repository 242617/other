package main

// #include <stdio.h>
// #include <stdlib.h>
import "C"
import "unsafe"

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
	cs := C.CString("hello!")
	defer C.free(unsafe.Pointer(cs))
	// C.fputs(cs, (*C.FILE)(C.stdout))
	return cs
}

//export SetString
func SetString(str string) *C.char {
	cs := C.CString(str)
	defer C.free(unsafe.Pointer(cs))
	// C.fputs(cs, (*C.FILE)(C.stdout))
	return cs
}
