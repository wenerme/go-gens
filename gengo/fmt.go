package gengo

import (
	"strings"

	"mvdan.cc/gofumpt/format"

	"github.com/wenerme/go-gens/gen"
	"golang.org/x/tools/imports"
)

func gofmt(f []byte) (i []byte, err error) {
	return format.Source(f, format.Options{
		ExtraRules: true,
	})
}

func Format(f *gen.File) (err error) {
	if !strings.HasSuffix(f.Name, ".go") {
		return
	}
	f.Content, err = gofmt(f.Content)
	if err == nil {
		f.Content, err = imports.Process(f.Name, f.Content, nil)
	}
	return
}
