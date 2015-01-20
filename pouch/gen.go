package main

import (
	"text/template"

	"github.com/ttacon/builder"
	"github.com/ttacon/pouch/pouch/defs"
)

func generateStructCode(toGen []*defs.StructInfo) ([]byte, error) {
	var (
		fileBytes      = builder.NewBuilder(nil)
		templateBuffer = builder.NewBuilder(nil)
	)
	for _, s := range toGen {
		err := structTmplt.Execute(templateBuffer, s)
		if err != nil {
			return nil, err
		}

		fileBytes.Write(templateBuffer.Bytes())
		templateBuffer.Reset()
	}

	return fileBytes.Bytes(), nil
}

func generateFunctions(toGen []*defs.StructInfo) ([]byte, error) {
	var (
		fileBytes = builder.NewBuilder(nil)
	)
	templateToGoThrough := []*template.Template{
		identifiableT,
		insertableT,
		tableablT,
		findableT,
		gettableT,
	}
	for _, s := range toGen {
		for _, templ := range templateToGoThrough {
			err := templ.Execute(fileBytes, s)
			if err != nil {
				return nil, err
			}
		}
	}
	return fileBytes.Bytes(), nil
}
