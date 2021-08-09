package gen

import "context"

type Template struct {
	Name           string // template name
	Filename       string // generated file name
	FilenameFormat func(ctx context.Context) string
	Skip           func(ctx context.Context) bool
}

func (t Template) Template() Template {
	return t
}
