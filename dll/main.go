package main

import "C"
import (
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const LibPath = "lib.dll"

func main() {
	log.SetFlags(log.Lshortfile)
	fmt.Println("start")

	GetInt()
	SetInt()
	GetString()
	SetString()

	fmt.Println("done")
}

func GetInt() {
	fmt.Println("GetInt")
	proto := syscall.NewLazyDLL(LibPath)
	getInt := proto.NewProc("GetInt")
	res, r2, err := getInt.Call()
	if r2 != 0 {
		panic(err)
	}
	fmt.Println(res)
}

func SetInt() {
	fmt.Println("SetInt")
	proto := syscall.NewLazyDLL(LibPath)
	setInt := proto.NewProc("SetInt")
	res, r2, err := setInt.Call(13)
	if r2 != 0 {
		panic(err)
	}
	fmt.Println(res)
}

func GetString() {
	fmt.Println("GetString")
	proto := syscall.NewLazyDLL(LibPath)
	getString := proto.NewProc("GetString")
	res, r2, err := getString.Call()
	if r2 != 0 {
		panic(err)
	}
	fmt.Println(string(res))
}

func SetString() {
	fmt.Println("SetString")
	proto := syscall.NewLazyDLL(LibPath)
	setString := proto.NewProc("SetString")
	str := "test string"
	prt := unsafe.Pointer(&str)
	uptr := (uintptr)(prt)
	res, r2, err := setString.Call(uptr)
	if r2 != 0 {
		panic(err)
	}
	fmt.Println(string(res))
}
