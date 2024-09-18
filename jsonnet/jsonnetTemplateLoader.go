package main

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
)

var fake = &jsonnet.NativeFunction{
	Name:   "fake",
	Params: ast.Identifiers{"x"},
	Func: func(x []interface{}) (interface{}, error) {
		bytes, err := json.Marshal(x[0])
		if err != nil {
			return nil, err
		}
		return CallGoFakeIt(strings.Trim(string(bytes), "\""))
	},
}

func Load(templateFilePath string, generateTestData bool) string {
	vm := jsonnet.MakeVM()

	vm.Importer(&jsonnet.FileImporter{})
	vm.ExtVar("generateTestData", strconv.FormatBool(generateTestData))
	vm.NativeFunction(fake)
	jsonStr, err := vm.EvaluateFile(templateFilePath)

	if err != nil {
		log.Panic(err.Error())
	}
	return jsonStr
}
