package tools

import (
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
)

func MakeDir(elem ...string) (p string, ok bool) {
	ok = true
	p = path.Join(elem...)
	_, err := os.Stat(p)
	if err == nil {
		return
	}
	if os.IsExist(err) {
		return
	}
	err = os.MkdirAll(p, os.ModePerm)
	if err != nil {
		ok = false
		return
	}
	return
}

func GetFunctionName(i interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	if size := len(fields); size > 0 {
		return fields[size-1]
	}
	return ""
}
