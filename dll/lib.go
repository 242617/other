// +build windows,amd64
package main

/*
struct GetStructs {
   char *id;
   char *hash;
};
*/
import "C"
import (
	"strings"
	"unsafe"
)

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

var getStructsSize = C.size_t(unsafe.Sizeof(C.struct_GetStructs{}))

//export GetStructs
func GetStructs(id, hash *C.char) unsafe.Pointer {
	ptr := C.malloc(getStructsSize)
	res := (*C.struct_GetStructs)(ptr)
	res.id, res.hash = hash, id
	return unsafe.Pointer(ptr)
}
