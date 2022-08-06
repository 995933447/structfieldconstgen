package main

import (
	"github.com/995933447/std-go/scan"
)

// structfieldconst . -findPkgPath ../wegod/internal/pkg/datamodels -outFile ../wegod/internal/pkg/db/enum/fields.go -prefix Field
func main() {
	findPkgPath := scan.OptStr("findPkgPath")
	outFile := scan.OptStr("outFile")
	structName := scan.OptStrDefault("struct", "")
	constPrefix := scan.OptStrDefault("prefix", "")
	constSuffix := scan.OptStrDefault("suffix", "")
	transStructFieldConstValFucName := scan.OptStrDefault("func", "snake")
	err := structFieldsToConsts(findPkgPath, structName, constPrefix, constSuffix, getTransFieldConstValFuncByName(transStructFieldConstValFucName), outFile)
	if err != nil {
		panic(err)
	}
}