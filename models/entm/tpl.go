package entm

import (
	"context"

	"github.com/wenerme/go-gens/gen"
)

type MetaModelTemplate struct {
	Name           string // template name.
	Filename       string
	FilenameFormat func(ctx context.Context, mm *EntityMetaModel) string
	Skip           func(ctx context.Context, mm *EntityMetaModel) bool
}

func (mt MetaModelTemplate) Template() gen.Template {
	t := gen.Template{
		Name:           mt.Name,
		Filename:       mt.Filename,
		FilenameFormat: nil,
		Skip: func(ctx context.Context) bool {
			mm, ok := ctx.Value(gen.ModelKey).(*EntityMetaModel)
			if !ok {
				return true
			}
			if mt.Skip != nil {
				return mt.Skip(ctx, mm)
			}
			return false
		},
	}
	if mt.FilenameFormat != nil {
		t.FilenameFormat = func(ctx context.Context) string {
			mm := ctx.Value(gen.ModelKey).(*EntityMetaModel)
			return mt.FilenameFormat(ctx, mm)
		}
	}
	return t
}
