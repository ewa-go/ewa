package egowebapi

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

type Operation struct {
	Description  string              `json:"description,omitempty"`
	Consumes     []string            `json:"consumes,omitempty"`
	Produces     []string            `json:"produces,omitempty"`
	Schemes      []string            `json:"schemes,omitempty"`
	Tags         []string            `json:"tags,omitempty"`
	Summary      string              `json:"summary,omitempty"`
	ExternalDocs *ExternalDocs       `json:"externalDocs,omitempty"`
	ID           string              `json:"operationId,omitempty"`
	Deprecated   bool                `json:"deprecated,omitempty"`
	Security     Security            `json:"security,omitempty"`
	Parameters   []*Parameter        `json:"parameters,omitempty"`
	Responses    map[string]Response `json:"responses,omitempty"`
}

type Parameter struct {
	Path             string  `json:"-"`
	Type             string  `json:"type,omitempty"`
	Description      string  `json:"description,omitempty"`
	Name             string  `json:"name,omitempty"`
	In               string  `json:"in,omitempty"`
	Required         bool    `json:"required,omitempty"`
	Schema           *Schema `json:"schema,omitempty"`
	CollectionFormat string  `json:"collectionFormat,omitempty"`
	AllowEmptyValue  bool    `json:"allowEmptyValue,omitempty"`
	Items            *Items  `json:"items,omitempty"`
}

type Schema struct {
	Ref   string `json:"$ref,omitempty"`
	Type  string `json:"type,omitempty"`
	Items *Items `json:"items,omitempty"`
	/*Discriminator string        `json:"discriminator,omitempty"`
	ReadOnly      bool          `json:"readOnly,omitempty"`
	XML           *XMLObject    `json:"xml,omitempty"`
	ExternalDocs  *ExternalDocs `json:"externalDocs,omitempty"`
	Example       interface{}   `json:"example,omitempty"`*/
}

/*type XMLObject struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Prefix    string `json:"prefix,omitempty"`
	Attribute bool   `json:"attribute,omitempty"`
	Wrapped   bool   `json:"wrapped,omitempty"`
}*/

type Response struct {
	Description string                 `json:"description"`
	Schema      *Schema                `json:"schema,omitempty"`
	Headers     map[string]Header      `json:"headers,omitempty"`
	Examples    map[string]interface{} `json:"examples,omitempty"`
}

type Header struct {
	Description string `json:"description,omitempty"`
	CommonValidations
	SimpleSchema
}

type CommonValidations struct {
	Maximum          *float64      `json:"maximum,omitempty"`
	ExclusiveMaximum bool          `json:"exclusiveMaximum,omitempty"`
	Minimum          *float64      `json:"minimum,omitempty"`
	ExclusiveMinimum bool          `json:"exclusiveMinimum,omitempty"`
	MaxLength        *int64        `json:"maxLength,omitempty"`
	MinLength        *int64        `json:"minLength,omitempty"`
	Pattern          string        `json:"pattern,omitempty"`
	MaxItems         *int64        `json:"maxItems,omitempty"`
	MinItems         *int64        `json:"minItems,omitempty"`
	UniqueItems      bool          `json:"uniqueItems,omitempty"`
	MultipleOf       *float64      `json:"multipleOf,omitempty"`
	Enum             []interface{} `json:"enum,omitempty"`
}

type SimpleSchema struct {
	Type             string      `json:"type,omitempty"`
	Nullable         bool        `json:"nullable,omitempty"`
	Format           string      `json:"format,omitempty"`
	Items            *Items      `json:"items,omitempty"`
	CollectionFormat string      `json:"collectionFormat,omitempty"`
	Default          interface{} `json:"default,omitempty"`
	Example          interface{} `json:"example,omitempty"`
}

type Items struct {
	Ref string `json:"$ref,omitempty"`
	CommonValidations
	SimpleSchema
}

// NewSchema Инициализация схемы для параметров
func NewSchema(i interface{}) *Schema {
	return &Schema{
		Ref: RefDefinition(i),
	}
}

// NewSchemaArray Инициализация схемы с массивом для параметров
func NewSchemaArray(i interface{}) *Schema {
	return &Schema{
		Type: TypeArray,
		Items: &Items{
			Ref: RefDefinition(i),
		},
	}
}

// RefDefinition Получаем имя модели, чтобы затем сформировать ссылку
func RefDefinition(i interface{}) string {

	var t reflect.Type
	value := reflect.ValueOf(i)
	if value.Type().Kind() == reflect.Ptr {
		t = reflect.Indirect(value).Type()
	} else {
		t = value.Type()
	}
	return t.Name()
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
func NewInFormData(name string, required bool, desc ...string) *Parameter {
	p := &Parameter{
		In:       InFormData,
		Name:     name,
		Required: required,
	}
	if desc != nil {
		p.Description = desc[0]
	}
	return p
}

// NewResponse Инициализация ответа
func NewResponse(schema *Schema, desc ...string) Response {
	r := Response{
		Schema:  schema,
		Headers: map[string]Header{},
	}
	if desc != nil {
		r.Description = desc[0]
	}
	return r
}

// AddHeader Добавить заголовок в ответ
func (r Response) AddHeader(name string, header Header) Response {
	r.Headers[name] = header
	return r
}

// NewHeader Инициализация заголовка
func NewHeader(t interface{}, nullable bool, desc ...string) Header {

	h := Header{
		SimpleSchema: SimpleSchema{
			Nullable: nullable,
		},
	}
	if desc != nil {
		h.Description = desc[0]
	}

	switch t.(type) {
	case string:
		h.Type = "string"
		break
	case int:
		h.Type = "integer"
		break
	case int8:
		h.Type = "integer"
		h.Format = "int8"
		break
	case int16:
		h.Type = "integer"
		h.Format = "int16"
		break
	case int32:
		h.Type = "integer"
		h.Format = "int32"
		break
	case int64:
		h.Type = "integer"
		h.Format = "int64"
		break
	case time.Time:
		h.Type = "string"
		h.Format = "date-time"
		break
	case bool:
		h.Type = "boolean"
		break
	}
	return h
}

// getParams Извлекаем пути из параметров
func (o Operation) getParams() (params []string) {
	for _, param := range o.Parameters {
		params = append(params, param.Path)
	}
	return
}

// addTag Добавить tag в список
func (o *Operation) addTag(tag string) {
	o.Tags = append(o.Tags, tag)
}
