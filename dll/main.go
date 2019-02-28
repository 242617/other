// +build windows,amd64
package main

import "C"
import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

var lib = syscall.NewLazyDLL("lib.dll")

func main() {
	log.SetFlags(log.Lshortfile)
	GetInt()
	SetInt()
	GetString()
	SetString()
}

func GetInt() {
	getInt := lib.NewProc("GetInt")
	res, _, err := getInt.Call()
	if err != syscall.Errno(0) {
		panic(err)
	}
	fmt.Println("GetInt:", res)
}

func SetInt() {
	setInt := lib.NewProc("SetInt")
	res, _, err := setInt.Call(13)
	if err != syscall.Errno(0) {
		panic(err)
	}
	fmt.Println("SetInt:", res)
}

func GetString() {
	getString := lib.NewProc("GetString")
	res, _, err := getString.Call()
	if err != syscall.Errno(0) {
		panic(err)
	}
	charPtr := (*C.char)(unsafe.Pointer(res))
	fmt.Println("GetString:", C.GoString(charPtr))
}

func SetString() {
	setString := lib.NewProc("SetString")
	cs := C.CString("Hey!")
	res, _, err := setString.Call(uintptr(unsafe.Pointer(cs)))
	if err != syscall.Errno(0) {
		panic(err)
	}
	charPtr := (*C.char)(unsafe.Pointer(res))
	fmt.Println("SetString:", C.GoString(charPtr))
}
