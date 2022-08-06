package main

import (
	"github.com/995933447/stringhelper-go"
)

var transStructFieldConstValFuncMap = map[string]transStructFieldConstValFunc{
	"snake": func(fieldName string) string {
		return stringhelper.Snake(fieldName)
	},
	"camel": func(fieldName string) string {
		return stringhelper.Camel(fieldName)
	},
}

func getTransFieldConstValFuncByName(name string) transStructFieldConstValFunc {
	if fn, ok := transStructFieldConstValFuncMap[name]; ok {
		return fn
	}
	panic("not support func:" + name)
}
