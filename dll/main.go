// +build windows,amd64
package main

/*

#include <stdlib.h>

struct GetStructs {
   char *id;
   char *hash;
};

*/
import "C"
import (
	"fmt"
	"syscall"
	"unsafe"
)

var lib = syscall.NewLazyDLL("lib.dll")

func main() {
	SetInt()
	GetInt()
	GetString()
	SetString()
	GetStructs()
}

func GetInt() {
	getInt := lib.NewProc("GetInt")
	res, _, err := getInt.Call()
	if err != syscall.Errno(0) {
		panic(err)
	}
	defer C.free(unsafe.Pointer(res))
	fmt.Println("GetInt:", res)
}

func SetInt() {
	setInt := lib.NewProc("SetInt")
	res, _, err := setInt.Call(13)
	if err != syscall.Errno(0) {
		panic(err)
	}
	defer C.free(unsafe.Pointer(res))
	fmt.Println("SetInt:", res)
}

func GetString() {
	getString := lib.NewProc("GetString")
	res, _, err := getString.Call()
	if err != syscall.Errno(0) {
		panic(err)
	}
	defer C.free(unsafe.Pointer(res))
	charPtr := (*C.char)(unsafe.Pointer(res))
	fmt.Println("GetString:", C.GoString(charPtr))
}

func SetString() {
	setString := lib.NewProc("SetString")

	cs := C.CString("Hey!")
	defer C.free(unsafe.Pointer(cs))

	res, _, err := setString.Call(uintptr(unsafe.Pointer(cs)))
	if err != syscall.Errno(0) {
		panic(err)
	}
	defer C.free(unsafe.Pointer(res))

	charPtr := (*C.char)(unsafe.Pointer(res))
	fmt.Println("SetString:", C.GoString(charPtr))
}

func GetStructs() {
	getStructs := lib.NewProc("GetStructs")

	id := C.CString("0x11e99157fc3feb")
	defer C.free(unsafe.Pointer(id))
	hash := C.CString("0x203ed75097c1d1")
	defer C.free(unsafe.Pointer(hash))

	res, _, err := getStructs.Call(uintptr(unsafe.Pointer(id)), uintptr(unsafe.Pointer(hash)))
	if err != syscall.Errno(0) {
		panic(err)
	}
	defer C.free(unsafe.Pointer(res))

	result := (*C.struct_GetStructs)(unsafe.Pointer(res))
	fmt.Println("GetStructs:", C.GoString(result.id), C.GoString(result.hash))
}
