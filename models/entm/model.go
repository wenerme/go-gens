package entm

import (
	"fmt"

	"github.com/huandu/xstrings"
	"github.com/jinzhu/inflection"
)

type EntityMetaModelField struct {
	Name    string
	NameZh  string `yaml:"nameZh"`
	Type    string
	SQLType string `yaml:"sqlType"`
	GoType  string `yaml:"goType"`
}

type EntityMetaModel struct {
	Name      string
	TableName string `yaml:"tableName"`
	Fields    []*EntityMetaModelField
}

func Normalize(mm *EntityMetaModel) error {
	if mm.TableName == "" {
		mm.TableName = inflection.Plural(xstrings.ToSnakeCase(mm.Name))
	}
	for _, f := range mm.Fields {
		if f.Type == "" {
			f.Type = "string"
		}
		if f.SQLType == "" {
			switch f.Type {
			case "string":
				f.SQLType = "text"
			default:
				f.SQLType = f.Type
			}
		}
		if f.GoType == "" {
			f.GoType = f.Type
		}
		if f.GoType == "float" {
			f.GoType = "float64"
		}
		if f.GoType == "" {
			return fmt.Errorf("no go type %q", f.Name)
		}
	}
	return nil
}
