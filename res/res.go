package res

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	raylib "github.com/gen2brain/raylib-go/raylib"
)

/*
#cgo windows CFLAGS: -I../raylib/src
#cgo windows LDFLAGS: -L../raylib/src
#include "res.h"
#include "stdlib.h"
*/
import "C"

type resource struct {
	data interface{}
	err  error
}

var (
	resources map[string]resource = map[string]resource{}
)

const tmp = ".tmp/"

//go:embed */**
var resfs embed.FS

func Init() error {
	err := LoadDir(".")
	if err != nil {
		return err
	}

	return nil
}

func Load() error {
	err := ProcFs()
	if err != nil {
		return err
	}
	return nil
}

func WriteFs() error {
	for path, res := range resources {
		if res.err == nil {
			err := os.MkdirAll(filepath.Dir(tmp+path), os.ModePerm)
			if err != nil {
				return err
			}

			err = os.WriteFile(tmp+path, res.data.([]byte), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CleanFs() {
	os.RemoveAll(tmp)
}

func LoadDir(cwd string) error {
	dir, err := resfs.ReadDir(cwd)
	if err != nil {
		return err
	}
	for _, file := range dir {
		path := file.Name()
		if cwd != "." {
			path = cwd + "/" + path
		}
		if file.IsDir() {
			LoadDir(path)
		} else {
			data, err := resfs.ReadFile(path)
			fmt.Printf("%+v: %+v\n", path, len(data))
			resources[path] = resource{data, err}
		}
	}
	return nil
}

func ProcFs() error {
	for path, res := range resources {
		if res.err == nil && strings.Contains(path, ".obj") {
			data, err := LoadObjFS(path)
			resources[path] = resource{data, err}
		}
	}

	return nil
}

func GetRes(path string) (interface{}, error) {
	if res, ok := resources[path]; ok {
		return res.data, res.err
	} else {
		return nil, fmt.Errorf("Resource not found: %s\n", path)
	}
}

//export GetData
func GetData(cfile *C.char, cdir *C.char) *C.char {
	file := C.GoString(cfile)
	dir := C.GoString(cdir)
	path := dir + "/" + file

	pos := strings.Index(path, "./")
	if pos == 0 {
		path = path[2:]
	}

	data, err := GetRes(path)
	if err == nil {
		return C.CString(string(data.([]byte)))
	}

	return nil
}

func LoadObjFS(path string) (interface{}, error) {
	res, ok := resources[path]
	if !ok {
		return nil, fmt.Errorf("Resource not found\n")
	}
	if res.err != nil {
		return nil, res.err
	}

	cpath := C.CString(path)
	cdata := C.CString(string(res.data.([]byte)))
	obj := C.LoadOBJ(cpath, cdata)
	C.free(unsafe.Pointer(cpath))
	C.free(unsafe.Pointer(cdata))

	mdl := *(*raylib.Model)(unsafe.Pointer(obj))
	return mdl, nil
}

func LoadIqmFS(path string) (interface{}, error) {
	res, ok := resources[path]
	if !ok {
		return nil, fmt.Errorf("Resource not found\n")
	}
	if res.err != nil {
		return nil, res.err
	}

	cpath := C.CString(path)
	cdata := C.CString(string(res.data.([]byte)))
	obj := C.LoadIQM(cpath, cdata)
	C.free(unsafe.Pointer(cpath))
	C.free(unsafe.Pointer(cdata))

	mdl := *(*raylib.Model)(unsafe.Pointer(obj))
	return mdl, nil
}

func LoadIqmAnimFS(path string) (interface{}, error) {
	res, ok := resources[path]
	if !ok {
		return nil, fmt.Errorf("Resource not found\n")
	}
	if res.err != nil {
		return nil, res.err
	}

	cpath := C.CString(path)
	cdata := C.CString(string(res.data.([]byte)))
	count := C.int(0)
	iqm := C.LoadAnimIQM(cpath, cdata, &count)
	C.free(unsafe.Pointer(cpath))
	C.free(unsafe.Pointer(cdata))

	anim := (*[1 << 24]raylib.ModelAnimation)(unsafe.Pointer(iqm))[:int(count)]
	return anim, nil
}
