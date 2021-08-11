package gen_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/wenerme/go-gens/gengqls"

	"github.com/Masterminds/sprig"
	"github.com/wenerme/go-gens/gen"
	"github.com/wenerme/go-gens/gengo"
	"github.com/wenerme/go-gens/models/entm"

	"github.com/stretchr/testify/assert"
)

func trimTxt(f *gen.File) error {
	if strings.HasSuffix(f.Name, ".txt") {
		f.Content = bytes.TrimSpace(f.Content)
	}
	return nil
}

func TestGen(t *testing.T) {
	tpl := template.New("test")
	tpl.Funcs(sprig.TxtFuncMap())
	template.Must(tpl.Parse(`
{{define "hello"}}
hello {{.}} !
{{end}}

{{define "go"}}
package gen
type {{.Name}} struct{
{{range $_,$f:= .Fields }}
{{$f.Name}} {{$f.GoType}}
{{end}}
}
{{end}}

{{define "gqls"}}
type {{.Name}} {
  {{- range $k, $v := .Fields}}
  """{{$v.NameZh}}"""
  {{$v.Name}}: {{$v.Type | title}}
  {{- end}}
}
{{end}}
`))

	load := false
	g := &gen.Generator{
		Debug:     true,
		Template:  tpl,
		Formatter: gen.Formatters(trimTxt, gengo.Format, gengqls.Format),
		Templates: []gen.IsTemplate{
			gen.Template{
				Name: "hello",
				FilenameFormat: func(ctx context.Context) string {
					return fmt.Sprintf("hello-%v.txt", ctx.Value(gen.ModelKey).(string))
				},
				Skip: func(ctx context.Context) bool {
					_, ok := ctx.Value(gen.ModelKey).(string)
					return !ok
				},
			},
			gen.Template{
				Name: "skip",
				Skip: func(ctx context.Context) bool {
					return true
				},
			},
			entm.MetaModelTemplate{
				Filename: "model.go",
				Name:     "go",
				Skip: func(ctx context.Context, mm *entm.EntityMetaModel) bool {
					return false
				},
			},
			entm.MetaModelTemplate{
				Filename: "model.graphqls",
				Name:     "gqls",
				Skip: func(ctx context.Context, mm *entm.EntityMetaModel) bool {
					return false
				},
			},
		},
		Loader: func() error {
			load = true
			return nil
		},
	}
	mm := &entm.EntityMetaModel{
		Name: "People",
		Fields: []*entm.EntityMetaModelField{
			{Name: "name"},
		},
	}
	assert.NoError(t, entm.Normalize(mm))
	assert.NoError(t, g.Generate("wener", mm))

	assert.True(t, load)
	assert.Equal(t, &gen.File{
		Name:    "hello-wener.txt",
		Content: []byte(`hello wener !`),
	}, g.Files[0])
	assert.Equal(t, &gen.File{
		Name: "model.go",
		Content: []byte(`package gen

type People struct {
	name string
}
`),
	}, g.Files[1])

	assert.NoError(t, g.Write(gen.WriteConfig{
		DryRun: true,
	}))
}
