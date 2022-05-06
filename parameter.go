package egowebapi

import (
	"regexp"
	"strings"
)

type Parameter struct {
	Path             string  `json:"-"`
	Type             string  `json:"type,omitempty"`
	Format           string  `json:"format,omitempty"`
	Description      string  `json:"description,omitempty"`
	Name             string  `json:"name,omitempty"`
	In               string  `json:"in,omitempty"`
	Required         bool    `json:"required,omitempty"`
	Schema           *Schema `json:"schema,omitempty"`
	CollectionFormat string  `json:"collectionFormat,omitempty"`
	AllowEmptyValue  bool    `json:"allowEmptyValue,omitempty"`
	Items            *Items  `json:"items,omitempty"`
}

// NewInPath Инициализация параметра in: path
func NewInPath(path string, required bool, desc ...string) *Parameter {

	// Извлекаем параметр из пути
	matches := regexp.MustCompile(`{(\w+)}`).FindStringSubmatch(path)
	if len(matches) < 2 {
		return nil
	}

	p := &Parameter{
		Path:     path,
		Name:     matches[1],
		In:       InPath,
		Required: required,
		Type:     TypeString,
	}
	if desc != nil {
		p.Description = desc[0]
	}

	return p
}

// NewInBody Инициализация параметра in: body
func NewInBody(required bool, schema *Schema, desc ...string) *Parameter {
	p := &Parameter{
		In:       InBody,
		Name:     InBody,
		Required: required,
		Schema:   schema,
		Type:     TypeString,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewInHeader Инициализация параметра in: header
func NewInHeader(name string, required bool, desc ...string) *Parameter {
	p := &Parameter{
		In:       InHeader,
		Name:     name,
		Type:     TypeString,
		Required: required,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewInQuery Инициализация параметра in: query
func NewInQuery(name string, required bool, desc ...string) *Parameter {
	p := &Parameter{
		In:       InQuery,
		Name:     name,
		Type:     TypeString,
		Required: required,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewInQueryArray Инициализация параметра in: query с типом массив
func NewInQueryArray(name, array string, required bool, desc ...string) *Parameter {
	var (
		enum        []interface{}
		defaultItem string
	)
	for i, a := range strings.Split(array, ",") {
		a = strings.Trim(a, " ")
		if i == 0 {
			defaultItem = a
		}
		enum = append(enum, a)
	}
	p := &Parameter{
		In:               InQuery,
		Name:             name,
		Type:             TypeArray,
		CollectionFormat: CollectionFormatMulti,
		Required:         required,
		Items: &Items{
			CommonValidations: CommonValidations{
				Enum: enum,
			},
			SimpleSchema: SimpleSchema{
				Type: TypeString,
				//CollectionFormat: CollectionFormatMulti,
				Default: defaultItem,
			},
		},
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewInFormData Инициализация параметра in: formData
func NewInFormData(name, t string, required bool, desc ...string) *Parameter {
	p := &Parameter{
		Name:     name,
		In:       InFormData,
		Required: required,
		Type:     t,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewParameter Инициализация параметра
func NewParameter(name string) *Parameter {
	return &Parameter{
		Name: name,
	}
}

// SetType Установка типа данных параметра
func (p *Parameter) SetType(t string) *Parameter {
	p.Type = t
	return p
}

// SetIn Установка типа параметра
func (p *Parameter) SetIn(i string) *Parameter {
	p.In = i
	return p
}

// SetFormat Установка формата данных параметра
func (p *Parameter) SetFormat(format string) *Parameter {
	p.Format = format
	return p
}

// SetRequired Установка флага обязательности параметра
func (p *Parameter) SetRequired(required bool) *Parameter {
	p.Required = required
	return p
}

// SetDescription Установка описания параметра
func (p *Parameter) SetDescription(desc string) *Parameter {
	p.Description = desc
	return p
}

// SetSchema Установка описания параметра
func (p *Parameter) SetSchema(schema *Schema) *Parameter {
	p.Schema = schema
	return p
}

// SetTypeFormat Установка описания параметра
func (p *Parameter) SetTypeFormat(t interface{}) *Parameter {
	p.Type, p.Format = setTypeFormat(t)
	return p
}
