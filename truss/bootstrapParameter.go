package truss

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// BootstrapParameter - struct to handle types and case conversions
type BootstrapParameter struct {
	Type  string
	Value string

	PascalCase string
	CamelCase  string
	KebabCase  string
	SnakeCase  string
	FlatCase   string
}

// NewBootstrapParameter - create a bootstrap parameter with the type string
func NewBootstrapParameter(value string) *BootstrapParameter {
	param := &BootstrapParameter{Type: "string", Value: value}
	param.generateCases()

	return param
}

func (c *BootstrapParameter) String() string {
	return c.Value
}

func (c *BootstrapParameter) generateCases() {
	c.PascalCase = strcase.ToCamel(c.Value)
	c.CamelCase = strcase.ToLowerCamel(c.Value)
	c.KebabCase = strcase.ToKebab(c.Value)
	c.SnakeCase = strcase.ToSnake(c.Value)
	c.FlatCase = strings.ReplaceAll(strcase.ToSnake(c.Value), "_", "")
}

// NewBootstrapParameterBool - create a bootstrap parameter with the type bool
func NewBootstrapParameterBool(value bool) *BootstrapParameter {
	var v string
	if value {
		v = "true"
	} else {
		v = "false"
	}

	param := &BootstrapParameter{Type: "bool", Value: v}
	param.generateCases()

	return param
}
