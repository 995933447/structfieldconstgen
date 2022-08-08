package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type (
	transStructFieldConstValFunc func(string) string

	FieldForTransConst struct {
		Name string
		ConstVal string
	}

	TransFieldsToConstsData struct {
		ConstPrefix string
		ConstSuffix string
		Fields []*FieldForTransConst
		OutPkg string
	}
)

func structFieldsToConsts(findPkgPath, structName string, constPrefix, constSuffix string, transConstValHandler transStructFieldConstValFunc, outFile string) error {
	fieldNames, err := findFieldsForTransConst(findPkgPath, structName)
	if err != nil {
		return err
	}

	outDir := filepath.Dir(outFile)
	if err = os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	transFieldsToConstsData := &TransFieldsToConstsData{
		ConstSuffix: constSuffix,
		ConstPrefix: constPrefix,
		OutPkg: filepath.Base(outDir),
	}
	for _, fieldName := range fieldNames {
		transFieldsToConstsData.Fields = append(transFieldsToConstsData.Fields, &FieldForTransConst{
			Name: fieldName,
			ConstVal: transConstValHandler(fieldName),
		})
	}

	tmpl := template.New("structFieldsToConsts")
	tmpl, err = tmpl.Parse(structFieldConstTmpl)
	if err != nil {
		return err
	}
	fp, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 7555)
	if err != nil {
		return err
	}
	if err = tmpl.Execute(fp, transFieldsToConstsData); err != nil {
		return err
	}

	return nil
}

func findFieldsForTransConst(findPkgPath string, structName string) ([]string, error) {
	pkgMap, err := parser.ParseDir(token.NewFileSet(), findPkgPath, func(info fs.FileInfo) bool {
		return !strings.Contains(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		return nil, err
	}

	var (
		existFieldNameMap = make(map[string]struct{})
		readyFieldNames []string
	)
	for _, pkg := range pkgMap {
		for _, file := range pkg.Files {
			for _, declNode := range file.Decls {
				genDeclNode, ok := declNode.(*ast.GenDecl)
				if !ok {
					continue
				}

				if genDeclNode.Tok != token.TYPE {
					continue
				}

				var typeSpec *ast.TypeSpec
				for _, spec := range genDeclNode.Specs {
					if typeSpec, ok = spec.(*ast.TypeSpec); ok {
						break
					}
				}

				if typeSpec == nil {
					continue
				}

				if structName != "" && typeSpec.Name.String() != structName {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				if structType.Fields.List == nil {
					continue
				}

				for _, field := range structType.Fields.List {
					if field.Names == nil {
						continue
					}
					for _, fieldIdent := range field.Names {
						if _, ok = existFieldNameMap[fieldIdent.Name]; ok {
							continue
						}
						existFieldNameMap[fieldIdent.Name] = struct{}{}
						readyFieldNames = append(readyFieldNames, fieldIdent.Name)
					}
				}
			}
		}
	}

	return readyFieldNames, nil
}