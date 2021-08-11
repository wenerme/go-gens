package gengqls

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/wenerme/go-gens/gen"
)

func Format(f *gen.File) error {
	if !strings.HasSuffix(f.Name, ".graphqls") {
		return nil
	}
	o := bytes.Buffer{}
	// no validate
	s, err := parser.ParseSchema(&ast.Source{
		Name:  f.Name,
		Input: string(f.Content),
	})
	if err != nil {
		lines := strings.Split(string(f.Content), "\n")
		c := strings.Builder{}
		for i, v := range lines {
			c.WriteString(fmt.Sprintf("%v: ", i))
			c.WriteString(v)
			c.WriteRune('\n')
		}
		log.Println(f.Name, "\n", c.String())
		return err
	}
	formatter.NewFormatter(&o).FormatSchemaDocument(s)
	f.Content = o.Bytes()
	return nil
}
